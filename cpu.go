package main

import log "github.com/sirupsen/logrus"

// CPU ..
type CPU struct {
	// Used to store memory addresses
	// only lowest (rightmost) 12 bits are used
	// since the capacity of RAM is 4k
	registerI uint16

	// delay and sound timer
	delay, sound uint8

	// used to store the current executing addr
	programCounter uint16

	// used to point to the topmost level of the stack
	stackPointer byte

	// registers
	register [16]byte

	// chip-8 allowing upto 16 levels of nested subroutines
	stack [16]uint16
}

// StepTimers : Update timer values per second according to the frequency of their clocks
func (cpu *CPU) StepTimers() {

	cpu.delay = MaxOf(0, cpu.delay-1)
	cpu.sound = MaxOf(0, cpu.sound-1)
}

func newCPU() *CPU {

	log.Info("Initing CPU..")

	cpu := &CPU{
		programCounter: ProgramAreaStart}
	cpu.delay = 0

	return cpu
}
