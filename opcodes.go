package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/mobile/event/key"
)

// Contains CHIP-8 instruction set of 36 instructions

// All instructions are 2 bytes long and are stored
// most-significant-byte first (Big endian). In memory, the first byte
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

// 00E0 - CLS
// Clear display
func (vm *VM) cls() {
	scr := vm.screen
	log.Debug("Clearing display")
	scr.clearDisplay()
	BufferToScreen(scr)

	vm.IncrementPC()
}

// 00EE - RET
// Return from sub-routine
func (vm *VM) ret() {

	cpu := vm.cpu
	log.Debugf("Returning from sub-routine: PC: %d and SP: %d", cpu.programCounter, cpu.stackPointer)
	cpu.stackPointer--
	cpu.programCounter = cpu.stack[cpu.stackPointer]

	vm.IncrementPC()
}

// 1nnn - JP addr
// Jump to location nnn
func (vm *VM) jp(nnn uint16) {
	// @discuss: should we validate the addr before setting it?
	cpu := vm.cpu
	log.Debugf("JMP to addr %d and PC: %d and SP: %d", nnn, cpu.programCounter, cpu.stackPointer)
	cpu.programCounter = nnn
}

// 2nnn - CALL addr
// Puts the current PC on the top of the stack. The PC is then
// set to nnn.
func (vm *VM) call(nnn uint16) {
	// should we validate the addr before setting it
	cpu := vm.cpu

	log.Debugf("CALL %d and PC: %d and SP: %d", nnn, cpu.programCounter, cpu.stackPointer)

	cpu.stack[cpu.stackPointer] = cpu.programCounter
	cpu.stackPointer++

	cpu.programCounter = nnn
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
func (vm *VM) se(x uint8, kk byte) {

	cpu := vm.cpu

	log.Debugf("SKIP NXT INS if Vx == kk, Vx: %d and kk: %d", cpu.register[x], kk)

	if cpu.register[x] == kk {
		vm.SkipInstruction()
	} else {
		vm.IncrementPC()
	}

}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
func (vm *VM) se_not(x uint8, kk byte) {

	cpu := vm.cpu

	log.Debugf("SKIP NXT INS if Vx != kk, Vx: %d and kk: %d", cpu.register[x], kk)

	if cpu.register[x] != kk {
		vm.SkipInstruction()
	} else {
		vm.IncrementPC()
	}
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
func (vm *VM) se_reg(x, y uint8) {
	cpu := vm.cpu

	log.Debugf("SKIP NXT INS if Vx == Vy, Vx: %d and kk: %d", cpu.register[x], cpu.register[y])

	if cpu.register[x] == cpu.register[y] {

		vm.SkipInstruction()
	} else {
		vm.IncrementPC()
	}
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
func (vm *VM) ld(vx uint8, data byte) {
	cpu := vm.cpu
	log.Debugf("LD byte: %d to Vx: %d", data, vx)
	cpu.register[vx] = data

	vm.IncrementPC()
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
func (vm *VM) add(vx uint8, data byte) {
	cpu := vm.cpu
	// log.Debugf("ADD-ING byte: %d to Vx: %d", data, cpu.register[vx])
	cpu.register[vx] += data

	vm.IncrementPC()
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
func (vm *VM) ld_reg(vx, vy uint8) {
	cpu := vm.cpu
	log.Debugf("LD data from Vy: %d to Vx: %d", vy, vx)
	cpu.register[vx] = cpu.register[vy]

	vm.IncrementPC()
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
func (vm *VM) or(vx, vy uint8) {
	cpu := vm.cpu

	cpu.register[vx] |= cpu.register[vy]

	vm.IncrementPC()
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
func (vm *VM) and(vx, vy uint8) {
	cpu := vm.cpu

	cpu.register[vx] &= cpu.register[vy]

	vm.IncrementPC()
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
func (vm *VM) xor(vx, vy uint8) {
	cpu := vm.cpu

	cpu.register[vx] ^= cpu.register[vy]

	vm.IncrementPC()
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
// The values of Vx and Vy are added together. If the result is greater than 8
// bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the
// result are kept, and stored in Vx.
func (vm *VM) add_reg(vx, vy uint8) {
	cpu := vm.cpu
	tmp := uint16(cpu.register[vx]) + uint16(cpu.register[vy])

	MAX := uint16(255)

	if tmp > MAX {
		cpu.register[0xF] = 1
	} else {
		cpu.register[0xF] = 0
	}

	// 8bits are kept(?!) or modulo 256 happens
	cpu.register[vx] += cpu.register[vy]

	vm.IncrementPC()
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
//
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx,
// and the results stored in Vx.
// @verify @check
func (vm *VM) sub_reg(vx, vy uint8) {
	cpu := vm.cpu

	if cpu.register[vx] >= cpu.register[vy] {
		cpu.register[0xF] = 1
	} else {
		cpu.register[0xF] = 0
	}

	cpu.register[vx] -= cpu.register[vy]

	vm.IncrementPC()
}

// 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
// @test
func (vm *VM) shr(vx, vy uint8) {
	cpu := vm.cpu

	cpu.register[0xF] = cpu.register[vx] & 1
	cpu.register[vx] = cpu.register[vx] >> 1

	vm.IncrementPC()
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.

// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func (vm *VM) subn(x, y uint8) {
	cpu := vm.cpu

	if cpu.register[y] > cpu.register[x] {
		cpu.register[0xF] = 1
	} else {
		cpu.register[0xF] = 0
	}

	cpu.register[x] = cpu.register[y] - cpu.register[x]

	vm.IncrementPC()
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0.
// Store the value of register VY shifted left one bit in register VX
func (vm *VM) shl(vx, vy uint8) {
	cpu := vm.cpu

	// set VF to MSB of VX
	cpu.register[0xF] = cpu.register[vx] >> 7
	// x = y << 1 OR x = x << 1: both works
	// check 8xyE notes at https://massung.github.io/CHIP-8/
	cpu.register[vx] = cpu.register[vx] << 1

	vm.IncrementPC()
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
func (vm *VM) sne(vx, vy uint8) {
	cpu := vm.cpu

	if cpu.register[vx] != cpu.register[vy] {
		vm.SkipInstruction()
	} else {
		vm.IncrementPC()
	}
}

// Annn - LD I, addr
// Set I = nnn.
func (vm *VM) ld_i(addr uint16) {
	cpu := vm.cpu
	log.Debugf("LD: Loading addr %d into register I", addr)
	cpu.registerI = addr

	vm.IncrementPC()
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
func (vm *VM) jp_add(addr uint16) {
	cpu := vm.cpu

	cpu.programCounter = addr + uint16(cpu.register[0])

	vm.IncrementPC()
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
func (vm *VM) rnd(vx uint8, kk byte) {
	cpu := vm.cpu

	// @improve: can we improve the rand here
	// also, need to @test this
	cpu.register[vx] = byte(RandInRange(0, 256)%0xFF) & kk

	vm.IncrementPC()
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.

// The interpreter reads n bytes from memory, starting at the address stored in I. These bytes are then displayed as sprites on screen at coordinates (Vx, Vy). Sprites are XORed onto the existing screen. If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen. See instruction 8xy3 for more information on XOR, and section 2.4, Display, for more information on the Chip-8 screen and sprites.
func (vm *VM) drw(vx, vy uint8, n uint8) {

	cpu := vm.cpu
	memory := vm.memory

	x := cpu.register[vx]
	y := cpu.register[vy]
	height := n

	log.Debugf("Drawing sprite at x: %d, y:%d", x, y)

	// @refactor: Take a slice of underlying memory instead?
	buf := make([]byte, n)
	startAddr := cpu.registerI

	// read N byte sprite data into buf starting from startAddr
	for i := uint16(0); i < uint16(height); i++ {
		buf[i] = memory.ram[startAddr+i]
	}

	scr := vm.screen

	// reset collision register
	cpu.register[0xF] = 0

	// display and update collision flag
	// j for height of the buffer
	for j := uint8(0); j < height; j++ {

		// spread each byte as 8 bits @test
		for i := uint8(0); i < 8; i++ {

			res := int((buf[j] >> i) & 1)

			// wrap around if required
			yLine := (y + j) % EmuHeight
			xLine := (x + (8 - i - 1)) % EmuWidth

			if scr.display[yLine][xLine] == 1 && res == 0 {
				cpu.register[0xF] = 1
			}

			scr.display[yLine][xLine] ^= res
		}
	}

	BufferToScreen(scr)

	vm.IncrementPC()
}

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
func (vm *VM) skp(vx uint8) {
	cpu := vm.cpu
	k := vm.keyboard

	vxData := cpu.register[vx]
	keyState := k.GetKeyState(vxData)

	log.Printf("Invoking ExA1: vxData: %d and mapped key: %s", vxData, k.reverseKeyboardMap[vxData])

	if keyState == key.DirPress {
		vm.SkipInstruction()
	} else {
		vm.IncrementPC()
	}
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
func (vm *VM) sknp(vx uint8) {
	cpu := vm.cpu
	k := vm.keyboard

	vxData := cpu.register[vx]
	keyState := k.GetKeyState(vxData)

	log.Printf("Invoking ExA1: vxData: %d and mapped key: %s", vxData, k.reverseKeyboardMap[vxData])

	if keyState != key.DirPress {
		vm.SkipInstruction()
	} else {
		vm.IncrementPC()
	}
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
func (vm *VM) ld_dt_in_vx(vx uint8) {
	cpu := vm.cpu

	cpu.register[vx] = cpu.delay

	vm.IncrementPC()
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (vm *VM) ld_key(vx uint8) {
	currTime := time.Now()

	k := vm.keyboard

	var pressedKeyCode key.Code

	log.Debugf("BEFORE: LastPressedKeyTime: %t and currTime: %t", k.lastPressedKey.time, currTime)

	// Blocked until we get a correct input
	for {
		if k.lastPressedKey.time.After(currTime) {
			pressedKeyCode = k.lastPressedKey.code
			break
		}
		// add 1 cycle sleep b/w polling
		time.Sleep(CPUTickerSpeed)
		// TODO: in future, can this be an event?
	}

	log.Debugf("AFTER: LastPressedKeyTime: %t and currTime: %t", k.lastPressedKey.time, currTime)

	cpu := vm.cpu

	val, _ := k.keyboardMap[pressedKeyCode]

	// reset pressed key
	k.keyboardState[pressedKeyCode] = key.DirRelease

	cpu.register[vx] = val

	vm.IncrementPC()
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
func (vm *VM) ld_dt(vx uint8) {
	cpu := vm.cpu

	cpu.delay = cpu.register[vx]

	vm.IncrementPC()
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
func (vm *VM) ld_st(vx uint8) {
	cpu := vm.cpu

	cpu.sound = cpu.register[vx]

	vm.IncrementPC()
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
// The values of I and Vx are added, and the results are stored in I.
func (vm *VM) add_i(vx uint8) {
	cpu := vm.cpu

	cpu.registerI = uint16(cpu.register[vx]) + cpu.registerI

	vm.IncrementPC()
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx. See section 2.4, Display, for more information on the Chip-8 hexadecimal font.
func (vm *VM) ld_font(x uint8) {
	cpu := vm.cpu

	digit := cpu.register[x]
	// offset of 5 bytes per digit, refer to chip8Fontset in memory.go
	cpu.registerI = uint16(DigitSpriteDataStart) + uint16(0x5*digit)

	vm.IncrementPC()
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.

// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
// @supershaddy impl, not sure if this works or not
func (vm *VM) bcd_ld(x uint8) {
	cpu := vm.cpu
	vxData := cpu.register[x]

	memory := vm.memory

	I := cpu.registerI

	memory.ram[I] = vxData / 100
	memory.ram[I+1] = (vxData / 10) % 10
	memory.ram[I+2] = vxData % 10

	vm.IncrementPC()
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
func (vm *VM) ld_i_to_vx(vx uint8) {
	cpu := vm.cpu
	memory := vm.memory

	for reg := uint8(0); reg <= vx; reg++ {
		// reading each byte into the register
		memory.ram[cpu.registerI+uint16(reg)] = cpu.register[reg]
	}

	vm.IncrementPC()
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
func (vm *VM) ld_vx(vx uint8) {
	cpu := vm.cpu
	memory := vm.memory
	addr := cpu.registerI

	for i := uint16(0); i <= uint16(vx); i++ {
		// reading each byte into the register
		cpu.register[i] = memory.ram[addr+i]
	}

	// cpu.registerI += uint16(vx) + 1

	vm.IncrementPC()
}

// IncrementPC makes PC point to next instruction
func (vm *VM) IncrementPC() {
	cpu := vm.cpu
	cpu.programCounter += uint16(2)
}

func (vm *VM) SkipInstruction() {
	// skipping two because the instruction is of 2
	// bytes size i.e. incrementing program counter by 2
	cpu := vm.cpu
	cpu.programCounter += uint16(4)
}

// Opcode ... @future management in this way will help
// us avoid ugly switch case, but I will probably
// implement this in next improvement
type Opcode struct {
	noParam [2]func()
}
