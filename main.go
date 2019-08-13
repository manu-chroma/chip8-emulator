package main

import (
	"log"
	"os"
	"time"
)

/// Refer to: http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
/// for good CHIP-8 reference mannual

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Booting up CHIP-8...")

	// parse the rom file
	args := os.Args

	// pick the last
	romFilePath := args[len(args)-1]
	log.Printf("Provided rom filepath: %s", romFilePath)

	conf := VMConfig{
		romFilePath: romFilePath}

	// create VM
	vm := InitVM(&conf)

	log.Println("Rom file: ", vm.memory.ram[ProgramAreaStart:ProgramAreaStart+vm.memory.romSize])

	for {
		vm.Tick()
		var t = (2 * 100 * time.Duration(1e6))
		time.Sleep(t)
	}
}
