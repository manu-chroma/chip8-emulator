package main

// Screen ..
// we draw graphics on screen through the use of sprites.
// A sprite is a group of bytes which are a binary representation of the desired picture.
// Chip-8 sprites can be up to 15 bytes, for a possible sprite size of 8x15
type Screen struct {
	display [64][32]bool
}
