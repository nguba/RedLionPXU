package device

// Modbus defines the interface for Modbus communication
type Modbus interface {
	SetUnitId(id uint8) error
	ReadRegisters(address, quantity uint16) ([]uint16, error)
	SetRegister(address uint16, regs uint16) error
	Close() error
}

// PxuReader defines the interface for reading device data
type PxuReader interface {
	ReadStats() (*Stats, error)
	ReadInfo() (*Info, error)
	ReadProfile() error
	Close() error
}

type PxuWriter interface {
	WriteSp(value uint16) error
}
