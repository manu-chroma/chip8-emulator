package main

import (
	"encoding/binary"
	"log"

	"golang.org/x/mobile/event/key"
)

// VM ...
type VM struct {
	cpu    *CPU
	screen *Screen
	memory *Memory
	// keyboardEvents propagation is needed since we read them from display screen
	// and require to access them in some of the instruction opcodes
	// this channel serves as a buffer for this
	// @verify todo: no deadlock condition should be there in case of empty buffer
	keyboardEvents chan key.Event
}

// VMConfig ...
type VMConfig struct {
	romFilePath string
}

// InitVM ...
func InitVM(vmConfig *VMConfig) *VM {

	vm := new(VM)

	vm.cpu = newCPU()
	vm.memory = newMemory()

	// todo: error handling
	vm.memory.LoadRomFile(vmConfig.romFilePath)

	// setup display
	vm.keyboardEvents = make(chan key.Event, 100)
	vm.screen = NewDisplay(vm.keyboardEvents)

	// so as to improve the abstraction
	// we can pass the actualKeyboard events
	// to keypad method, which will transform
	// to the hex keyboard of chip-8 and pass
	// to another transformed/actual hex keyboard channel
	// @discuss
	// keypad channel

	return vm
}

// ReadOpcode checks the memory and the current state of cpu
// and returns the current opcode which is a 2 byte in size
func (vm *VM) ReadOpcode() (uint16, error) {

	memory := vm.memory
	cpu := vm.cpu

	// pick out program counter
	pc := cpu.programCounter
	// no need to off-set since we read the rom
	// directly at the start of memory.ram buffer
	// instead of 0x0200
	// pc += 0x200

	// read two bytes of data and concat
	// this approach works as well, but found a better
	// direct function call to do the same
	// op1 := memory.ram[pc]
	// op2 := memory.ram[pc+1]
	// concat the 2 bytes of code
	// opcode := (uint16(op1) << 8) | uint16(op2)

	opcode := binary.LittleEndian.Uint16(memory.ram[pc : pc+2])

	log.Printf("Identified opcode: %d", opcode)

	cpu.programCounter += 2

	return opcode, nil
	// return 0, errors.New("Failed to read")
}

// Tick ...
func (vm *VM) Tick() {

	opcode, err := vm.ReadOpcode()

	if err != nil {
		log.Fatal(err)
	}

	vm.executeOpcode(opcode)

	log.Printf("Executed: %d", opcode)
}

// This will be our massive switch statement for now
// but I intend to improve it with some matching and all
// to leverage function pointer array
func (vm *VM) executeOpcode(opcode uint16) {

	if opcode == 0 {
		// no operation
	}

	// Atanomy of a CHIP-8 opcode
	// 2 bytes opcode
	//   1st nibble 2nd nibble      3rd nibble 4th nibble
	// |_______________________|  |_______________________|
	//        upperByte                  lowerByte

	// terminology used for variables and opcode below
	// @todo
	// note: add here, same docs have been provided
	// in opcodes.go for easy understanding

	// todo: assign correct values for all here
	upperByte := opcode & 0xFF00
	// in most signifiant -> to least significant order
	firstNibble := upperByte & 0xF
	// secondNibble := 0
	thirdNibble := 0
	fourthNibble := 0

	mmm := opcode & 0xFFF
	// xy := opcode & 0xFF0

	var x, y uint8
	var kk byte

	lowerByte := opcode & 0x00FF

	// NOP
	if opcode == 0 {

	} else if upperByte == 0 {
		if lowerByte == 0xE0 {
			vm.cls()
		} else if lowerByte == 0xEE {
			vm.ret()
		}
	} else if firstNibble == 1 {
		// 1nnn
		vm.jp(mmm)
	} else if firstNibble == 2 {
		// 2nnn
		vm.call(mmm)
	} else if firstNibble == 3 {
		// 3xkk
		vm.se(x, kk)
	} else if firstNibble == 4 {
		// 4xkk
		vm.se_not(x, kk)
	} else if firstNibble == 5 {
		// 5xy0
		vm.se_reg(x, y)
	} else if firstNibble == 6 {
		// 6xkk
		vm.ld(x, kk)
	} else if firstNibble == 7 {
		// 7xkk
		vm.add(x, kk)
	} else if firstNibble == 8 {
		if fourthNibble == 0 {
			// 8xy0
			vm.ld_reg(x, y)
		} else if fourthNibble == 1 {
			// 8xy1
			vm.or(x, y)
		} else if fourthNibble == 2 {
			// 8xy2
			vm.and(x, y)
		} else if fourthNibble == 3 {
			// 8xy3
			vm.xor(x, y)
		} else if fourthNibble == 4 {
			// 8xy4
			vm.add_reg(x, y)
		} else if fourthNibble == 5 {
			// 8xy5
			vm.add_reg(x, y)
		}

		// is there as 6 as well?

	} else if firstNibble == 9 {
		// 9xy0
		vm.sne(x, y)
	} else if firstNibble == 0xA {
		// Annn
		vm.ld_i(mmm)
	} else if firstNibble == 0xB {
		// Bnnn
		vm.jp_add(mmm)
	} else if firstNibble == 0xC {
		// Cxkk
		vm.rnd(x, kk)
	} else if firstNibble == 0xD {
		// Dxyn
		n := uint8(fourthNibble)
		vm.drw(x, y, n)
	} else if firstNibble == 0xE {
		if thirdNibble == 0x9 {
			// Ex9E
			vm.skp(x)
		} else if thirdNibble == 0xA {
			// ExA1
			vm.sknp(x)
		}
	} else if firstNibble == 0xF {
		// last remaining in series
	}

}
