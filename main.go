package main

import (
	"log"
	"os"
	"time"
)

/// Refer to: http://mattmik.com/files/chip8/mastering/chip8.html

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

	// @hack: @fix sleep for 2 seconds to ensure the window is up and running
	time.Sleep(2 * 1000 * time.Duration(1e6))

	// todo: don't crash at this line if no rom is loaded 
	log.Println("Rom file: ", vm.memory.ram[ProgramAreaStart:ProgramAreaStart+vm.memory.romSize])

	// TODO: the screen should be running in the main go thread.
	// https://stackoverflow.com/a/57474359/1180321

	for {
		vm.Tick()
		var t = (2 * 100 * time.Duration(1e6))
		time.Sleep(t)
	}

}
