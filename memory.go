package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

// Memory util constants
const (
	RAMEndAddr             = 0xFFF  // 4095
	InterpBlackListAdrrEnd = 0x1FF  // 511
	RAMSize                = 0x1000 // 4096

	SpriteDataStart = 0x000
	SpriteDataEnd   = 0x1FF
)

// Memory module contains the RAM
type Memory struct {
	ram [RAMSize]byte
}

// InitMemory ...
func newMemory() *Memory {
	m := new(Memory)
	return m
}

// LoadRomFile ...
func (m *Memory) LoadRomFile(romFilePath string) {

	// verify valid, readable file

	f, err := os.Open(romFilePath) // Error handling elided for brevity.
	if err != nil {
		log.Printf("Not able to load the rom file: %v", err)
		// log.Fatal("not able to open rom file..", err)
		return
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f) // Error handling elided for brevity.

	if err != nil {
		log.Printf("Not able to read rom data into buffer")
	}

	// todo: where to copy the rom data?

	// src -> destiantion
	// todo: check if this can be improved
	// copy(m.romMemory[:], buf.Bytes()[:])

}
