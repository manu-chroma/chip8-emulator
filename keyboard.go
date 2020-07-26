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

// Keyboard, todo: document members
type Keyboard struct {
	keypad             [4][4]byte
	keyboardState      map[key.Code]key.Direction
	keyboardMap        map[key.Code]byte
	reverseKeyboardMap map[byte]key.Code
	lastPressedKey     LastPressedKey
}

func newKeyboard() *Keyboard {
	k := Keyboard{}

	k.keyboardState = make(map[key.Code]key.Direction)

	// define keypad
	keypad := [4][4]byte{
		{1, 2, 3, 0xC},
		{4, 5, 6, 0xD},
		{7, 8, 9, 0xE},
		{0xA, 0, 0xb, 0xF}}

	k.keypad = keypad
	k.keyboardMap = map[key.Code]byte{
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

	k.reverseKeyboardMap = make(map[byte]key.Code)
	for keyy, val := range k.keyboardMap {
		k.reverseKeyboardMap[val] = keyy
	}

	return &k
}

func eligibleKeyEvent(event key.Event, k *Keyboard) bool {
	_, present := k.keyboardMap[event.Code]

	// skip DirRelease and DirNone events
	// DirNone indicates that the key was not released
	// or Pressed; it remains static..

	return present && (event.Direction == key.DirPress || event.Direction == key.DirNone)
}

// ProcessKeyEvent ..
func (k *Keyboard) ProcessKeyEvent(event key.Event) {
	if !eligibleKeyEvent(event, k) {
		return
	}

	// record last key press
	k.lastPressedKey.code = event.Code
	k.lastPressedKey.time = time.Now()

	// Put the `Pressed` event inside the keyboard state.
	k.keyboardState[event.Code] = key.DirPress
}

// GetKeyState ..
func (k *Keyboard) GetKeyState(chip8Key byte) key.Direction {
	res, _ := k.reverseKeyboardMap[chip8Key]
	dir, _ := k.keyboardState[res]

	// side-effect: reset key state since we're consuming this
	k.keyboardState[res] = key.DirRelease

	return dir
}
