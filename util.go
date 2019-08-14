package main

import (
	"fmt"
	"math/rand"
)

// MinOf arbitary no. of bytes
func MinOf(vars ...byte) byte {
	mini := vars[0]
	for _, i := range vars {
		if mini > i {
			mini = i
		}
	}

	return mini
}

// MaxOf arbitary no. of bytes
func MaxOf(vars ...byte) byte {
	maxi := vars[0]
	for _, i := range vars {
		if i > maxi {
			maxi = i
		}
	}

	return maxi
}

// RandInRange returns an int between the range
// [min, max). Will panic if (max - min) <= 0
func RandInRange(min, max int) int {

	return rand.Intn(max-min) + min
}

// HexOf ...
func HexOf(num uint16) string {
	return fmt.Sprintf("%x", num)
}

// HexOfByte ...
func HexOfByte(num byte) string {
	return fmt.Sprintf("%x", num)
}

// Reverse a string
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
