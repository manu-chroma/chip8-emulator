package main

import (
	"golang.org/x/mobile/event/key"
	"time"
)

// LastPressedKey stores the what and when the key was pressed
type LastPressedKey struct {
	code key.Code
	time time.Time
}

var lastPressedKey LastPressedKey

var keypad = [4][4]byte{
	{1, 2, 3, 0xC},
	{4, 5, 6, 0xD},
	{7, 8, 9, 0xE},
	{0xA, 0, 0xb, 0xF}}

var keyboardMap = map[key.Code]byte{
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
	key.CodeF: keypad[3][3],
}

var reverseKeyboardMap map[byte]key.Code

func generateReverseKeyMap() {
	reverseKeyboardMap = make(map[byte]key.Code)
	for k, v := range keyboardMap {
		reverseKeyboardMap[v] = k
	}
}

func eligibleKeyEvent(event key.Event) bool {
	_, present := keyboardMap[event.Code]

	// skip DirRelease and DirNone events
	// DirNone indicates that the key was not released
	// or Pressed; it remains static..

	return present && (event.Direction == key.DirPress)
}

// ProcessKeyEvent ..
func ProcessKeyEvent(event key.Event) {
	if !eligibleKeyEvent(event) {
		return
	}

	// record last key press
	lastPressedKey.code = event.Code
	lastPressedKey.time = time.Now()

	// Put the `Pressed` event inside the keyboard state.
	keyboardState[event.Code] = key.DirPress
}

var keyboardState map[key.Code]key.Direction

// InitKeyboard ..
func InitKeyboard() {
	lastPressedKey = LastPressedKey{}
	keyboardState = make(map[key.Code]key.Direction)
	generateReverseKeyMap()
}

// GetKeyState ..
func GetKeyState(chip8Key byte) key.Direction {
	res, _ := reverseKeyboardMap[chip8Key]
	dir, _ := keyboardState[res]

	// side-effect: reset key state since we're consuming this
	keyboardState[res] = key.DirRelease

	return dir
}
