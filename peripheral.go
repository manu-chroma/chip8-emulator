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

// InitKeyboard ..
func (kb *Keyboard) InitKeyboard() {
	kb.keypad = [4][4]byte{
		{1, 2, 3, 0xC},
		{4, 5, 6, 0xD},
		{7, 8, 9, 0xE},
		{0xA, 0, 0xb, 0xF}}

	keyboardMap = map[key.Code]byte{
		key.Code0: kb.keypad[3][1],
		key.Code1: kb.keypad[0][0],
		key.Code2: kb.keypad[0][1],
		key.Code3: kb.keypad[0][2],
		key.Code4: kb.keypad[1][0],
		key.Code5: kb.keypad[1][1],
		key.Code6: kb.keypad[1][2],
		key.Code7: kb.keypad[2][0],
		key.Code8: kb.keypad[2][1],
		key.Code9: kb.keypad[2][2],
		key.CodeA: kb.keypad[3][0],
		key.CodeB: kb.keypad[3][2],
		key.CodeC: kb.keypad[0][3],
		key.CodeD: kb.keypad[1][3],
		key.CodeE: kb.keypad[2][3],
		key.CodeF: kb.keypad[3][3]}
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
