package main

// VM ...
type VM struct {
	cpu    *CPU
	screen *Screen
	memory *Memory
}

// InitVM ...
func InitVM() *VM {

	vm := new(VM)

	vm.cpu = newCPU()
	vm.memory = newMemory()
	vm.screen = new(Screen)

	return vm
}
