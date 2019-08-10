package main

import "fmt"

// Contains CHIP-8 instruction set of 36 instructions

// All instructions are 2 bytes long and are stored
// most-significant-byte first. In memory, the first byte
// of each instruction should be located at an even
// addresses. If a program includes sprite data, it should
// be padded so any instructions following it will be
// properly situated in RAM.

// In these listings, the following variables are used:

// nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
// n or nibble - A 4-bit value, the lowest 4 bits of the instruction
// x - A 4-bit value, the lower 4 bits of the high byte of the instruction
// y - A 4-bit value, the upper 4 bits of the low byte of the instruction
// kk or byte - An 8-bit value, the lowest 8 bits of the instruction

// Clear display
// 00E0 - CLS
func cls() {

}

// Return from sub-routine
// 00EE - RET
func ret() {

	fmt.Println("Returning from sub-routine")

}

//
// 1nnn - JP addr
func (vm *VM) one(nnn uint16) {
	// should we validate the addr before setting it
	vm.cpu.programCounter = nnn
}

// Opcode ..
type Opcode struct {
	noParam [2]func()
}

func (op *Opcode) initOpcodeSet() {
	op.noParam = [2]func(){cls, ret}
}
