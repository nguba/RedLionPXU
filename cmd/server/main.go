package main

import (
	"flag"
	"fmt"
	"github.com/nguba/RedLionPXU/internal/device"
	"github.com/nguba/RedLionPXU/public/api"
	"log"
	"net"
	"time"
)

var (
	unit = flag.Int("unit", 5, "Unit Id configured for the device")
	mock = flag.Bool("mock", false, "Use a mock modbus implementation when testing without the device")
)

// DefaultConfiguration returns a default configuration for COM3
func DefaultConfiguration() *device.Configuration {
	return &device.Configuration{
		URL:      "rtu://COM3",
		Speed:    device.DefaultSpeed,
		DataBits: 8,
		Parity:   "none",
		Timeout:  500 * time.Millisecond,
	}
}

func main() {

	flag.Parse()

	unitId := device.UnitId(*unit)

	var modbus device.Modbus
	var err error

	if *mock {
		log.Println("using mock modbus implementation to impersonate the device")
		modbus = device.NewMockModbus()
	} else {
		cfg := DefaultConfiguration()
		modbus, err = device.NewModbusDevice(cfg)
	}
	if err != nil {
		log.Fatal(err)
	}

	pxu, err := device.NewPxu(unitId, modbus, device.DefaultTimeout, device.DefaultRetries)
	defer func(pxu *device.Pxu) {
		_ = pxu.Close()
	}(pxu)

	port := 5000 + *unit
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer func(lis net.Listener) {
		_ = lis.Close()
	}(lis)

	server, err := api.NewServer(pxu, lis)
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
