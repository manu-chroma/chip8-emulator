package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"time"
)

// Refer to: http://mattmik.com/files/chip8/mastering/chip8.html
// Excellent guide to understanding everything about chip8 emulation
const (
	CPUTickerSpeed = time.Duration(5) * time.Millisecond
)

var (
	EmulatorTick = time.NewTicker(CPUTickerSpeed)
)

func main() {
	setupLogging()

	log.Info("Booting up CHIP-8...")

	conf := parseConfig()
	vm := InitVM(&conf)

	// todo: document
	go func() {
		// @hack: sleep for 2 seconds to ensure the window (in screen struct) is up and running
		time.Sleep(2 * time.Second)
		log.Debugln("\n\n Rom file: ",
			vm.memory.ram[ProgramAreaStart:ProgramAreaStart+vm.memory.romSize])

		for {
			select {
			case <-EmulatorTick.C:
				vm.Tick()
			}
		}
	}()

	// Running display in the main go routine due to:
	// https://stackoverflow.com/a/57474359/1180321
	// Throws hard error when running on different routine
	// on macOS
	vm.InitDisplay(vm.keyboard)
}

func parseConfig() VMConfig {
	// Read romFilePath from cmd args
	romFilePath := flag.String("rom", "", "Rom File to execute on the interpreter")
	flag.Parse()

	if *romFilePath == "" {
		log.Fatal("Rom file path missing..")
	}

	log.Info("Provided rom filepath: %s", *romFilePath)
	conf := VMConfig{
		romFilePath: *romFilePath}

	return conf
}
