package main

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

// 00E0 - CLS
// Clear display
func (vm *VM) cls() {
	scr := vm.screen
	scr.clearDisplay()
}

// 00EE - RET
// Return from sub-routine
func (vm *VM) ret() {

	cpu := vm.cpu

	// todo: throwing error
	cpu.programCounter = cpu.stack[cpu.stackPointer]
	cpu.stackPointer--
}

// 1nnn - JP addr
// Jump to location nnn
func (vm *VM) jp(nnn uint16) {
	// @discuss: should we validate the addr before setting it?

	cpu := vm.cpu

	cpu.programCounter = nnn
}

// 2nnn - CALL addr
// The interpreter increments the stack pointer, then puts
// the current PC on the top of the stack. The PC is then
// set to nnn.
func (vm *VM) call(nnn uint16) {
	// should we validate the addr before setting it

	cpu := vm.cpu

	cpu.stackPointer++
	cpu.stack[cpu.stackPointer] = cpu.programCounter
	cpu.programCounter = nnn
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
func (vm *VM) se(x uint8, kk byte) {

	cpu := vm.cpu

	if cpu.register[x] == kk {
		// skipping two because the instruction is of 2
		// bytes size
		cpu.programCounter += 2
	}

}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
func (vm *VM) se_not(x uint8, kk byte) {

	cpu := vm.cpu

	if cpu.register[x] != kk {
		// skipping two because the instruction is of 2
		// bytes size
		cpu.programCounter += 2
	}
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
func (vm *VM) se_reg(x, y uint8) {
	cpu := vm.cpu

	if cpu.register[x] == cpu.register[y] {
		// skipping two because the instruction is of 2
		// bytes size
		cpu.programCounter += 2
	}
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
func (vm *VM) ld(vx uint8, data byte) {
	cpu := vm.cpu

	cpu.register[vx] = data
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
func (vm *VM) add(vx uint8, data byte) {
	cpu := vm.cpu

	cpu.register[vx] += data
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
func (vm *VM) ld_reg(vx, vy uint8) {
	cpu := vm.cpu

	cpu.register[vx] = cpu.register[vy]
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
func (vm *VM) or(vx, vy uint8) {
	cpu := vm.cpu

	vxData := cpu.register[vx]
	vyData := cpu.register[vy]
	cpu.register[vx] = vxData | vyData
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
func (vm *VM) and(vx, vy uint8) {
	cpu := vm.cpu

	vxData := cpu.register[vx]
	vyData := cpu.register[vy]
	cpu.register[vx] = vxData & vyData
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
func (vm *VM) xor(vx, vy uint8) {
	cpu := vm.cpu

	vxData := cpu.register[vx]
	vyData := cpu.register[vy]
	cpu.register[vx] = vxData ^ vyData
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.

// TODO the other ones that are in the middle

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
func (vm *VM) sne(vx, vy uint8) {
	cpu := vm.cpu

	if cpu.register[vx] != cpu.register[vy] {
		cpu.programCounter += 2
	}
}

// Annn - LD I, addr
// Set I = nnn.
func (vm *VM) ld_i(addr uint16) {
	cpu := vm.cpu

	cpu.registerI = addr
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
func (vm *VM) jp_add(addr uint16) {
	cpu := vm.cpu

	cpu.programCounter = addr + uint16(cpu.register[0])
}

// todo

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
func (vm *VM) ld_dt_in_vx(vx uint8) {
	cpu := vm.cpu

	cpu.register[vx] = cpu.delay
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
// All execution stops until a key is pressed, then the value of that key is stored in Vx.

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
func (vm *VM) ld_dt(vx uint8) {
	cpu := vm.cpu

	cpu.delay = cpu.register[vx]
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
func (vm *VM) ld_st(vx uint8) {
	cpu := vm.cpu

	cpu.sound = cpu.register[vx]
}

// todo

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
func (vm *VM) ld_i_to_vx(vx uint8, addr uint16) {
	cpu := vm.cpu
	memory := vm.memory

	for reg := uint8(0); reg <= vx; reg++ {
		// reading each byte into the register
		memory.ram[addr] = cpu.register[reg]
	}

}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
func (vm *VM) ld_vx(vx uint8, addr uint16) {
	cpu := vm.cpu
	memory := vm.memory

	for reg := uint8(0); reg <= vx; reg++ {
		// reading each byte into the register
		cpu.register[reg] = memory.ram[addr]
	}
}

// Opcode ..
type Opcode struct {
	noParam [2]func()
}

// func (op *Opcode) initOpcodeSet() {
// 	op.noParam = [2]func(){(vm) cls, ret}
// }
