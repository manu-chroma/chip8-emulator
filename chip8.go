package main

import (
	"encoding/binary"
	"errors"
	"log"

	"golang.org/x/mobile/event/key"
)

// ErrOpcodeNotImplemented ...
var ErrOpcodeNotImplemented = errors.New("Opcode not implemented yet!")

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

	log.Printf("Curr PC val: %d", pc)

	// no need to off-set since we directly map
	// the rom data at 0x200 in the memory.ram buffer
	// and init the PC with 0x200 val

	// read two bytes of data and concat
	// (what's happening in the function call)
	// op1 := memory.ram[pc]
	// op2 := memory.ram[pc+1]
	// concat the 2 bytes of code
	// opcode := (uint16(op1) << 8) | uint16(op2)

	opcode := binary.BigEndian.Uint16(memory.ram[pc : pc+2])

	hexRep := HexOf(opcode)
	log.Printf("**** Identified opcode :: %s ****\n", hexRep)

	cpu.programCounter += 2

	return opcode, nil
}

// Tick ...
func (vm *VM) Tick() {

	opcode, err := vm.ReadOpcode()

	if err != nil {
		// figure out a way to more gracefully
		// end the program when ROM execution is
		// completed
		// or we can have a special opcode for that
		// but that would be messing with the spec
		log.Fatal(err)
	}

	cpu := vm.cpu

	log.Printf("Before executing Opcode: %s, PC: %d, SP: %d", HexOf(opcode), cpu.programCounter, cpu.stackPointer)
	log.Print(cpu.register)
	vm.executeOpcode(opcode)

	log.Printf("Executed: %s", HexOf(opcode))
}

// This will be our massive switch statement for now
// but I intend to improve it with some matching and all
// to leverage function pointer array
func (vm *VM) executeOpcode(opcode uint16) {

	if opcode == 0 {
		// no operation
		log.Printf("NO OP code called! %s", HexOf(opcode))
	}

	// Atanomy of a CHIP-8 opcode
	// 2 bytes opcode
	//   1st nibble 2nd nibble      3rd nibble 4th nibble
	// |_______________________|  |_______________________|
	//        upperByte                  lowerByte

	// terminology used for variables and opcode below
	// In these listings, the following variables are used:

	// nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
	// n or nibble - A 4-bit value, the lowest 4 bits of the instruction
	// x - A 4-bit value, the lower 4 bits of the high byte of the instruction
	// y - A 4-bit value, the upper 4 bits of the low byte of the instruction
	// kk or byte - An 8-bit value, the lowest 8 bits of the instruction

	// NOTE: same docs (above) have been provided
	// in opcodes.go for easy understanding

	upperByte := byte(opcode >> 8) // & 0xFF00
	lowerByte := byte(opcode & 0xFF)

	// In most signifiant -> to least significant order
	firstNibble := upperByte >> 4
	secondNibble := upperByte & 0xF
	thirdNibble := lowerByte >> 4
	fourthNibble := lowerByte & 0xF

	mmm := opcode & 0xFFF
	// xy := opcode & 0xFF0

	x := secondNibble
	y := thirdNibble
	kk := lowerByte

	log.Printf("UB: %s LB: %s", HexOfByte(upperByte), HexOfByte(lowerByte))

	log.Printf("FstN: %s SN: %s TN: %s FthN: %s", HexOfByte(firstNibble), HexOfByte(secondNibble), HexOfByte(thirdNibble), HexOfByte(fourthNibble))

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
		if lowerByte == 0x07 {
			// Fx07
			vm.ld_dt_in_vx(x)
		} else if lowerByte == 0x15 {
			// Fx15
			vm.ld_dt(x)
		} else if lowerByte == 0x18 {
			// Fx18
			vm.ld_st(x)
		} else if lowerByte == 0x1E {
			// Fx1E
			vm.add_i(x)
		} else if lowerByte == 0x29 {
			// Fx29
			vm.ld_font(x)
		} else if lowerByte == 0x33 {
			// Fx33
			vm.bcd_ld(x)
		} else if lowerByte == 0x55 {
			// Fx55
			vm.ld_i_to_vx(x)
		} else if lowerByte == 0x65 {
			// Fx65
			vm.ld_vx(x)
		}
	}
}
