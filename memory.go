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

	DigitSpriteDataStart = 0x000
	DigitSpriteDataEnd   = 0x1FF

	DisplayDataStart = 0x100
	DisplayDataEnd   = 0x1FF

	ProgramAreaStart = 0x200
	ProgramAreaEnd   = 0xFFF
)

var chip8Fontset = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80} // F

// Memory module contains the RAM
type Memory struct {
	ram     [RAMSize]byte
	romSize int
}

// InitMemory ...
func newMemory() *Memory {
	m := &Memory{}
	setDigitDataInRAM(m)
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
	_, err = io.Copy(buf, f)

	if err != nil {
		log.Printf("Not able to read rom data into buffer")
	}

	// todo: where to copy the rom data?

	// src -> destiantion
	// todo: check if this can be improved
	copy(m.ram[ProgramAreaStart:], buf.Bytes()[:])

	m.romSize = buf.Len() // expressed as num of bytes

	log.Printf("Rom buffer size is: %d", m.romSize)

	log.Println("Sucessfully copied rom file into ram buffer")
}

// copy
func setDigitDataInRAM(m *Memory) {

	startAddr := DigitSpriteDataStart
	for i := 0; i < 80; i++ {
		m.ram[startAddr] = chip8Fontset[i]
		startAddr++
	}

	log.Print("Completed copying of chip8Fontset in the RAM memory")
	log.Printf("Current value of startAddr is: %s", HexOf(uint16(startAddr)))
}
