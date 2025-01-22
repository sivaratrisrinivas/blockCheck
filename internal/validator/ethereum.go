package validator

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/sha3"
)

var (
	// Basic Ethereum address format: 0x followed by 40 hex characters
	addressRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

// IsValidAddress checks if the given string is a valid Ethereum address
func IsValidAddress(address string) bool {
	return addressRegex.MatchString(address)
}

// IsChecksumAddress checks if the address has valid EIP-55 mixed-case checksum
func IsChecksumAddress(address string) bool {
	if !IsValidAddress(address) {
		return false
	}

	// Convert the address to checksum format
	checksummed, _ := ToChecksumAddress(address)

	// Compare with the input address
	return address == checksummed
}

// ToChecksumAddress converts an Ethereum address to mixed-case checksum format
func ToChecksumAddress(address string) (string, error) {
	if !IsValidAddress(address) {
		return "", fmt.Errorf("invalid ethereum address format")
	}

	addr := strings.ToLower(address[2:])
	hash := Keccak256([]byte(addr))
	result := "0x"

	for i := 0; i < len(addr); i++ {
		// For each character in the address
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}

		if addr[i] >= '0' && addr[i] <= '9' {
			result += string(addr[i])
		} else {
			if hashByte > 7 {
				result += strings.ToUpper(string(addr[i]))
			} else {
				result += string(addr[i])
			}
		}
	}
	return result, nil
}

// Keccak256 calculates the Keccak-256 hash of a byte slice
func Keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}
