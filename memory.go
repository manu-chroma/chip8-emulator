package main

// Memory util constants
const (
	RAMEndAddr             = 0xFFF  // 4095
	InterpBlackListAdrrEnd = 0x1FF  // 511
	RAMSize                = 0x1000 // 4096
)

// Memory module contains the RAM
type Memory struct {
	ram [RAMSize]byte
}

// InitMemory ...
func (m *Memory) InitMemory() {
}
