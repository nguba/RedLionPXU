package device

// Modbus defines the interface for Modbus communication
type Modbus interface {
	SetUnitId(id UnitId) error
	ReadRegister(address uint16) (uint16, error)
	ReadRegisters(address, quantity uint16) ([]uint16, error)
	SetRegister(address uint16, regs uint16) error
	SetRegisters(startAddr uint16, values []uint16) error
	Close() error
}

type PidController interface {
	UpdateSetpoint(value float64) error
	Close() error
}
