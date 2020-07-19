package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
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
	f, err := os.Open(romFilePath)
	if err != nil {
		log.Fatalf("Not able to load the rom file: %v", err)
		return
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)

	if err != nil {
		log.Infof("Not able to read rom data into buffer")
	}

	// Directly map the rom data at 0x200 in the memory.ram buffer
	// and init the PC with 0x200 val
	// This emulates the way the actual implementation works..
	// The 0x0-0x1FF range is for actual CHIP-8 emulator logic.
	copy(m.ram[ProgramAreaStart:], buf.Bytes()[:])

	m.romSize = buf.Len() // expressed as num of bytes

	log.Infof("Rom buffer size is: %d", m.romSize)
	log.Infoln("Successfully copied rom file into ram buffer")
}

// copy font-set into RAM
func setDigitDataInRAM(m *Memory) {

	for i := 0; i < 80; i++ {
		m.ram[i] = chip8Fontset[i]
	}

	log.Info("Completed copying of chip8 Font-set in the RAM memory")
}
