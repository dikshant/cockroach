package macaddr

import "testing"

func BenchmarkParseMACOriginal(b *testing.B) {
	macStr := "01:23:45:67:89:ab"
	var dest MACAddr

	for i := 0; i < b.N; i++ {
		if err := ParseMAC(macStr, &dest); err != nil {
			b.Fatalf("ParseMACOriginal failed: %v", err)
		}
	}
}

func BenchmarkParseMACOptimized(b *testing.B) {
	macStr := "01:23:45:67:89:ab"
	var dest MACAddr

	for i := 0; i < b.N; i++ {
		if err := ParseMACFast(macStr, &dest); err != nil {
			b.Fatalf("ParseMACOptimized failed: %v", err)
		}
	}
}
