// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package macaddr

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/errors"
)

// MACAddr is the representation of a MAC address. The uint64 takes 8-bytes for
// both 6 and 8 byte Macaddr.
type MACAddr uint64

// MACAddrSize is the size of a 6 byte MACAddr.
const MACAddrSize = 6

// ParseINet parses postgres style MACAddr types. While Go's net.ParseMAC
// supports MAC adddresses upto 20 octets in length, postgres does not. This
// function returns an error if the MAC address is longer than 64 bits. See
// TestParseMAC for examples.
func ParseMAC(s string, dest *MACAddr) error {
	hwAddr, err := net.ParseMAC(s)
	if err != nil {
		return pgerror.WithCandidateCode(
			errors.Errorf("could not parse %q as macaddr. invalid MAC address", s),
			pgcode.InvalidTextRepresentation)
	}

	var macInt uint64
	switch len(hwAddr) {
	case 6, 8:
		for _, b := range hwAddr {
			macInt = (macInt << 8) | uint64(b)
		}
	default:
		return pgerror.WithCandidateCode(
			errors.Errorf("could not parse %q as macaddr. invalid MAC address length", len(hwAddr)),
			pgcode.NumericValueOutOfRange)
	}

	*dest = MACAddr(macInt)
	return nil
}

// hibits returns the higher 32 bits of the MAC address.
func hibits(a MACAddr) uint32 {
	return uint32(uint64(a) >> 32)
}

// lobits returns the lower 32 bits of the MAC address.
func lobits(a MACAddr) uint32 {
	return uint32(uint64(a) & 0xFFFFFFFF)
}

// CompareMACs compares two MAC addresses using their high and low bits.
func (m MACAddr) Compare(addr MACAddr) int {
	if hibits(m) < hibits(addr) {
		return -1
	} else if hibits(m) > hibits(addr) {
		return 1
	} else if lobits(m) < lobits(addr) {
		return -1
	} else if lobits(m) > lobits(addr) {
		return 1
	} else {
		return 0
	}
}

// MacAddrNot performs bitwise NOT on the MAC address.
func MacAddrNot(addr MACAddr) MACAddr {
	return MACAddr(^uint64(addr))
}

// MacAddrAnd performs bitwise AND between two MAC addresses.
func MacAddrAnd(addr1, addr2 MACAddr) MACAddr {
	return MACAddr(uint64(addr1) & uint64(addr2))
}

// MacAddrOr performs bitwise OR between two MAC addresses.
func MacAddrOr(addr1, addr2 MACAddr) MACAddr {
	return MACAddr(uint64(addr1) | uint64(addr2))
}

// String implements the stringer interface for MACAddr.
func (m MACAddr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		byte(m>>40),
		byte(m>>32),
		byte(m>>24),
		byte(m>>16),
		byte(m>>8),
		byte(m))
}

// RandMACAddr generates a random MACAddr.
func RandMACAddr(rng *rand.Rand) MACAddr {
	buf := make([]byte, 6)
	var mac net.HardwareAddr

	_, err := crand.Read(buf)
	if err != nil {
		panic(err)
	}

	// Set the local bit
	buf[0] |= 2

	mac = append(mac, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	var macAddr MACAddr
	for _, b := range mac {
		macAddr = (macAddr << 8) | MACAddr(b)
	}

	return macAddr
}

// ToBuffer appends the MACAddr encoding to a buffer and returns the final buffer.
func (m *MACAddr) ToBuffer(appendTo []byte) []byte {
	// Convert the uint64 MAC address to a 6-byte array.
	macBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(macBytes, uint64(*m))

	// Append only the last 6 bytes of the 8-byte array to the buffer.
	appendTo = append(appendTo, macBytes[2:]...)

	return appendTo
}

// FromBuffer populates an MACAddr with data from a byte slice, returning the
// remaining buffer or an error.
func (m *MACAddr) FromBuffer(data []byte) ([]byte, error) {
	if len(data) < 6 {
		return nil, errors.AssertionFailedf("MACAddr decoding error: insufficient data, got %d bytes", len(data))
	}
	// Create a slice to hold the 6-byte MAC address.
	macBytes := data[:6]

	// Update macAddr to point to the new data.
	*m = MACAddr(binary.BigEndian.Uint64(append([]byte{0, 0}, macBytes...)))

	return data[6:], nil
}
