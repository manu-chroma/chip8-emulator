package main

import (
	"golang.org/x/mobile/event/key"
)

// VM ...
type VM struct {
	cpu    *CPU
	screen *Screen
	memory *Memory
	// mouseEvents propagation is needed since we read them from display screen
	// and require to access them in some of the instruction opcodes
	// this channel serves as a buffer for this
	// @verify todo: no deadlock condition should be there in case of empty buffer
	mouseEvents chan key.Event
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

	// todo: error
	vm.memory.LoadRomFile(vmConfig.romFilePath)

	// setup display
	vm.mouseEvents = make(chan key.Event, 100)
	vm.screen = NewDisplay(vm.mouseEvents)

	// for x := range mouseEvents {
	// 	log.Printf("Reciveed: %s", x.Code)
	// }

	// TODO
	// and keypad channel

	return vm
}
