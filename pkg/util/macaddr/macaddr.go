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
	"net"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/errors"
)

// MacAddr is the representation of the Mac address. The uint64 takes 8-bytes for
// both 6 and 8 byte Macaddr.
type Addr uint64

// ParseINet parses postgres style MACAddr types. While Go's net.ParseMAC
// supports MAC adddresses upto 20 octets in length, postgres does not. This
// function returns an error if the MAC address is longer than 64 bits. See
// TestParseMAC for examples.
func ParseMAC(s string, dest *Addr) error {
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

	*dest = Addr(macInt)
	return nil
}

// hibits returns the higher 32 bits of the MAC address.
func hibits(a Addr) uint32 {
	return uint32(uint64(a) >> 32)
}

// lobits returns the lower 32 bits of the MAC address.
func lobits(a Addr) uint32 {
	return uint32(uint64(a) & 0xFFFFFFFF)
}

// CompareMACs compares two MAC addresses using their high and low bits.
func Compare(a1, a2 Addr) int {
	if hibits(a1) < hibits(a2) {
		return -1
	} else if hibits(a1) > hibits(a2) {
		return 1
	} else if lobits(a1) < lobits(a2) {
		return -1
	} else if lobits(a1) > lobits(a2) {
		return 1
	} else {
		return 0
	}
}

// MacAddrNot performs bitwise NOT on the MAC address.
func MacAddrNot(addr Addr) Addr {
	return Addr(^uint64(addr))
}

// MacAddrAnd performs bitwise AND between two MAC addresses.
func MacAddrAnd(addr1, addr2 Addr) Addr {
	return Addr(uint64(addr1) & uint64(addr2))
}

// MacAddrOr performs bitwise OR between two MAC addresses.
func MacAddrOr(addr1, addr2 Addr) Addr {
	return Addr(uint64(addr1) | uint64(addr2))
}
