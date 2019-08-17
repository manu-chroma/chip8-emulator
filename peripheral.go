package main

import (
	"errors"

	"golang.org/x/mobile/event/key"
)

// Keyboard ..
type Keyboard struct {
	keypad [4][4]byte
}

var keyboardMap map[key.Code]byte

var keypad = [4][4]byte{
	{1, 2, 3, 0xC},
	{4, 5, 6, 0xD},
	{7, 8, 9, 0xE},
	{0xA, 0, 0xb, 0xF}}

// InitKeyboard ..
func InitKeyboard() {

	keyboardMap = map[key.Code]byte{
		key.Code0: keypad[3][1],
		key.Code1: keypad[0][0],
		key.Code2: keypad[0][1],
		key.Code3: keypad[0][2],
		key.Code4: keypad[1][0],
		key.Code5: keypad[1][1],
		key.Code6: keypad[1][2],
		key.Code7: keypad[2][0],
		key.Code8: keypad[2][1],
		key.Code9: keypad[2][2],
		key.CodeA: keypad[3][0],
		key.CodeB: keypad[3][2],
		key.CodeC: keypad[0][3],
		key.CodeD: keypad[1][3],
		key.CodeE: keypad[2][3],
		key.CodeF: keypad[3][3]}
}

// Chip8Key ..
// we need to have method here to capture the normal keyboard keyevent and return the mapping to the CHIP-8 keyboard
func Chip8Key(normalKey key.Code) (byte, error) {

	val, prs := keyboardMap[normalKey]

	if prs == true {
		return val, nil
	}

	return 0, errors.New("key not available")
}
