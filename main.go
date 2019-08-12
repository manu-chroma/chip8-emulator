package main

import (
	"log"
	"os"
)

/// Refer to: http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
/// for good CHIP-8 reference mannual

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Booting up CHIP-8...")

	// parse the rom file
	args := os.Args

	// pick the las
	romFilePath := args[len(args)-1]
	log.Printf("Provided rom filepath: %s", romFilePath)

	conf := VMConfig{
		romFilePath: romFilePath}

	// create VM
	vm := InitVM(&conf)

	log.Println("Rom file: ", vm.memory.ram[:vm.memory.romSize])

	vm.Tick()

	// for {
	// }

	// start processing
	// for {
	// 	vm.cpu.Tick()
	// }

	// experimenting with function pointer for opcodes
	// var op Opcode
	// op.initOpcodeSet()

}
