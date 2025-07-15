package main

import (
	"flag"
	"fmt"
	device2 "github.com/nguba/RedLionPXU/internal/device"
	"log"
	"time"
)

func main() {
	var (
		unitID = flag.Uint("unit", 6, "Modbus unit ID (default: 6)")
		port   = flag.String("port", "COM3", "Serial port (default: COM3)")
		infoF  = flag.Bool("info", false, "Print device information")
		statsF = flag.Bool("stats", false, "Print device statistics")
		profF  = flag.Bool("profile", false, "Read the profile")
	)

	flag.Parse()

	// Create configuration
	cfg := &device2.Configuration{
		URL:      fmt.Sprintf("rtu://%s", *port),
		Speed:    38400,
		DataBits: 8,
		Parity:   "none",
		Timeout:  500 * time.Millisecond,
	}

	// Create real client
	client, err := device2.NewModbusHandler(cfg)
	if err != nil {
		log.Fatalf("Failed to instantiate modbus handler: %v", err)
	}

	// Create PXU instance
	pxu, err := device2.NewPxu(uint8(*unitID), client, 3*time.Second, 30)
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}
	defer func(pxu *device2.Pxu) {
		err := pxu.Close()
		if err != nil {
			log.Fatalf("Failed to close controller: %v", err)
		}
	}(pxu)

	if infoF != nil && *infoF {
		info, err := pxu.ReadInfo()
		if err != nil {
			log.Fatalf("Failed to read info: %v", err)
		}
		fmt.Println(info)
	}

	if statsF != nil && *statsF {
		showStats(pxu)
	}

	if profF != nil && *profF {
		for i := uint16(0); i < 16; i++ {
			profile, err := pxu.ReadProfile(i)
			if err != nil {
				log.Fatalf("Failed to read profile: %v", err)
			}
			fmt.Println(profile)
		}
	}

	val := 35.0
	if err := pxu.UpdateSetpoint(val); err != nil {
		log.Fatalf("Failed to write Sp: %v", err)
	}

	if err := pxu.Run(); err != nil {
		log.Fatalf("Failed to start controller: %v", err)
	}
	showStats(pxu)
	time.Sleep(time.Second * 3)

	if err := pxu.Stop(); err != nil {
		log.Fatalf("Failed to stop controller: %v", err)
	}
	showStats(pxu)

}

func showStats(pxu *device2.Pxu) {
	stats, err := pxu.ReadStats()
	if err != nil {
		log.Fatalf("Failed to read stats: %v", err)
	}
	fmt.Println(stats)
}
