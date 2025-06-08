package main

import (
	"RedLionPXU/internal/device"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	var (
		unitID = flag.Uint("unit", 6, "Modbus unit ID (default: 6)")
		port   = flag.String("port", "COM3", "Serial port (default: COM3)")
	)
	flag.Parse()

	// Create configuration
	cfg := &device.Configuration{
		URL:      fmt.Sprintf("rtu://%s", *port),
		Speed:    38400,
		DataBits: 8,
		Parity:   "none",
		Timeout:  500 * time.Millisecond,
	}

	// Create real client
	client, err := device.NewPxuClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create PXU instance
	pxu, err := device.NewPxu(uint8(*unitID), client, 5*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to create PXU: %v", err)
	}
	defer pxu.Close()

	stats, err := pxu.ReadStats()
	if err != nil {
		log.Fatalf("Failed to read stats: %v", err)
	}
	fmt.Println(stats)

	err = pxu.ReadProfile(0)
	if err != nil {
		log.Fatalf("Failed to read profile: %v", err)
	}
}
