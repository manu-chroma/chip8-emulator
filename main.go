package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// Refer to: http://mattmik.com/files/chip8/mastering/chip8.html
// Excellent guide to understanding everything about chip8 emulation
const (
	CPUTickerSpeed = time.Duration(20) * time.Millisecond
)

func main() {

	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.InfoLevel)


	log.Info("Booting up CHIP-8...")

	// pick the last
	romFilePath := os.Args[len(os.Args)-1]
	log.Info("Provided rom filepath: %s", romFilePath)

	conf := VMConfig{
		romFilePath: romFilePath}

	// create VM
	vm := InitVM(&conf)
	emulatorTick := time.NewTicker(CPUTickerSpeed)

	go func() {

		// @hack: sleep for 2 seconds to ensure the window (in screen struct) is up and running
		time.Sleep(2 * time.Second)

		// todo: don't crash at this line if no rom is loaded
		log.Infoln()
		log.Infoln("Rom file: ", vm.memory.ram[ProgramAreaStart:ProgramAreaStart+vm.memory.romSize])

		for {
			select {
			case <-emulatorTick.C:
				// startT := time.Now()
				vm.Tick()
				// endT := time.Now()

				// todo: something seems wrong here ...
				// log.("Start time: %s, end time: %s", startT.String(), endT.String())
				// log.Info("Time of execution: %s", fmtDuration(time.Since(startT)))
			}
		}

	}()

	// TODO: the screen should be running in the main go thread.
	// https://stackoverflow.com/a/57474359/1180321
	// throws hard error when running the code on macOS
	vm.InitDisplay()

}
