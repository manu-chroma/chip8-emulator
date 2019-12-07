package main

import (
	"log"
	"os"
	"time"
)

// Refer to: http://mattmik.com/files/chip8/mastering/chip8.html
// Excellent guide to understanding everything about chip8 emulation

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Booting up CHIP-8...")

	// pick the last
	romFilePath := os.Args[len(os.Args)-1]
	log.Printf("Provided rom filepath: %s", romFilePath)

	conf := VMConfig{
		romFilePath: romFilePath}

	// create VM
	vm := InitVM(&conf)
	emulatorTick := time.NewTicker(time.Duration(20) * time.Millisecond)

	go func() {

		// @hack: sleep for 2 seconds to ensure the window (in screen struct) is up and running
		time.Sleep(2 * time.Second)

		// todo: don't crash at this line if no rom is loaded
		log.Println()
		log.Println("Rom file: ", vm.memory.ram[ProgramAreaStart:ProgramAreaStart+vm.memory.romSize])

		for {
			select {
			case <-emulatorTick.C:
				startT := time.Now()
				vm.Tick()
				endT := time.Now()

				// todo: something seems wrong here ...
				log.Printf("Start time: %s, end time: %s", startT.String(), endT.String())
				log.Printf("Time of execution: %d", time.Since(startT).Round(time.Nanosecond))
			}
		}

	}()

	// TODO: the screen should be running in the main go thread.
	// https://stackoverflow.com/a/57474359/1180321
	// throws hard error when running the code on macOS
	vm.InitDisplay()

}
