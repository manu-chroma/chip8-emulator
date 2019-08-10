package main

// we draw graphics on screen through the use of sprites.
// A sprite is a group of bytes which are a binary representation of the desired picture.
// Chip-8 sprites can be up to 15 bytes, for a possible sprite size of 8x15
const (
	Row = 32
	Col = 64
)

// Screen ...
type Screen struct {
	display [Row][Col]bool
}

func (scr *Screen) clearDisplay() {

	for i := 0; i < Row; i++ {
		for j := 0; j < Col; j++ {
			scr.display[i][j] = false
		}
	}
}
