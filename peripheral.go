package main

// Keyboard ..
type Keyboard struct {
	keypad [4][4]uint8
}

// InitKeyboard ..
func (kb *Keyboard) InitKeyboard() {
	kb.keypad = [4][4]uint8{
		{1, 2, 3, 0xC},
		{4, 5, 6, 0xD},
		{7, 8, 9, 0xE},
		{0xA, 0, 0xb, 0xF}}

}

// Chip8Key ..
// we need to have method here to capture the normal keyboard keyevent and return the mapping to the CHIP-8 keyboard
func Chip8Key(normalKey int) (int, error) {
	// TODO

	return -1, nil
}
