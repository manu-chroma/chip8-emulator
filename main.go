package main

import "log"

/// Refer to: http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
/// for good CHIP-8 reference mannual

func main() {
	log.Println("Booting up CHIP-8...")

	// parse the rom file
	/*
		romFilePath := flag.String("rom", "default", "path to rom file")
		flag.Parse()

		// romFilePath = ""

		if *romFilePath == "" {
			log.Fatalf("Rom path not provided! Exiting...")
		}

		log.Printf("Provided rom filepath: %s", *romFilePath)

	*/

	// experimenting with function pointer for opcodes
	// var op Opcode
	// op.initOpcodeSet()

	// create VM
	_ = newCPU()

}
