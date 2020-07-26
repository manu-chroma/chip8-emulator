package main

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	"time"
)

// VM contains the whole state of emulator
type VM struct {
	cpu      *CPU
	screen   *Screen
	memory   *Memory
	keyboard *Keyboard
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
	vm.keyboard = newKeyboard()

	vm.memory.LoadRomFile(vmConfig.romFilePath)

	// Setup counting timers
	// timers should be tweaked at 60 Hz,
	// roughly every 16 millisecond
	// TODO: derive from CPU speed in some way?
	ticker := time.NewTicker(time.Duration(16666) * time.Microsecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				vm.cpu.StepTimers()
			}
		}
	}()

	return vm
}

// InitDisplay ..
func (vm *VM) InitDisplay(keyboard *Keyboard) {
	vm.NewDisplay(keyboard)
}

// ReadOpcode checks the memory and the current state of cpu
// and returns the current opcode which is a 2 byte in size
func (vm *VM) ReadOpcode() (uint16, error) {

	memory := vm.memory
	cpu := vm.cpu

	pc := cpu.programCounter

	// Read two bytes of data and concat
	// (what's happening in the below function call)
	// op1 := memory.ram[pc]
	// op2 := memory.ram[pc+1]
	// concat the 2 bytes of code
	// opcode := (uint16(op1) << 8) | uint16(op2)
	opcode := binary.BigEndian.Uint16(memory.ram[pc : pc+2])

	return opcode, nil
}

// Tick executes one OPCODE at a time
func (vm *VM) Tick() {

	opcode, _ := vm.ReadOpcode()
	cpu := vm.cpu

	log.Debug(cpu.register)
	vm.executeOpcode(opcode)
}

// Our massive switch statement
// @future: Can this be improved with rust like matching and
// how can function pointer array be leveraged to make this better
func (vm *VM) executeOpcode(opcode uint16) {

	// Notes to self while understanding the opcode anatomy and execution
	// Anatomy of a CHIP-8 opcode
	// Length of every opcode: 2 bytes
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

	// NOTE: Above docs have been provided
	// in opcodes.go for easy reference

	upperByte := byte(opcode >> 8) // & 0xFF00
	lowerByte := byte(opcode & 0xFF)

	// In most significant -> to least significant order
	firstNibble := upperByte >> 4
	secondNibble := upperByte & 0xF
	thirdNibble := lowerByte >> 4
	fourthNibble := lowerByte & 0xF

	mmm := opcode & 0xFFF

	x := secondNibble
	y := thirdNibble
	kk := lowerByte

	if opcode == 0 {
		// NOP
		log.Infof("NO OP code called! %s", HexOf(opcode))
	} else if upperByte == 0 {
		if lowerByte == 0xE0 {
			vm.cls()
		} else if lowerByte == 0xEE {
			vm.ret()
		} else {
			log.Info("Execute machine language subroutine at address NNN")
			log.Fatal("Opcode not implemented...")
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
			vm.sub_reg(x, y)
		} else if fourthNibble == 6 {
			// 8xy6
			vm.shr(x, y)
		} else if fourthNibble == 7 {
			// 8xy7
			vm.subn(x, y)
		} else if fourthNibble == 0xE {
			// 8xyE
			vm.shl(x, y)
		}
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
		} else if lowerByte == 0x0A {
			// Fx0A
			vm.ld_key(x)
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
