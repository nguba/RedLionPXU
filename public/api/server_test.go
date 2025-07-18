package api

import (
	"context"
	"github.com/nguba/RedLionPXU/internal/device"
	v2 "github.com/nguba/RedLionPXU/public/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

const (
	bufSize = 1024 * 1024
	unit    = device.UnitId(5)
)

var (
	lis    *bufconn.Listener
	modbus *device.MockModbus
)

// initialize a listener that doesn't depend on a networked implementation
func init() {
	lis = bufconn.Listen(bufSize)
	modbus = device.NewMockModbus()
}

func setupTestServer(t *testing.T) v2.RedLionPxuClient {
	t.Helper() // Mark this as a helper function

	// setup server
	lis := bufconn.Listen(1024 * 1024)
	pxu, err := device.NewPxu(unit, modbus, time.Second, 3)

	svc, err := NewServer(pxu, lis)
	t.Cleanup(func() {
		svc.Stop()
	})

	if err := svc.Start(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// setup client
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	return v2.NewRedLionPxuClient(conn)
}

func TestApi_GetStats(t *testing.T) {
	client := setupTestServer(t)

	err := modbus.SetRegisters(0, modbus.GetStatsRegister())
	if err != nil {
		t.Fatalf("failed to set registers: %v", err)
	}
	t.Cleanup(func() {
		modbus.Reset()
	})

	got, err := client.GetStats(context.Background(), &v2.GetStatsRequest{})
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	if got.Stats.Pv != 25.5 {
		t.Errorf("GetStats returned wrong pv: %v", got.Stats.Pv)
	}
	if got.Stats.Sp != 30.4 {
		t.Errorf("GetStats returned wrong sp: %v", got.Stats.Sp)
	}
	if got.Stats.Out1 != true {
		t.Errorf("GetStats returned wrong out1: %v", got.Stats.Out1)
	}
	if got.Stats.Out2 != false {
		t.Errorf("GetStats returned wrong out2: %v", got.Stats.Out2)
	}
	if got.Stats.Rs != device.Run.String() {
		t.Errorf("GetStats returned wrong rs: %v", got.Stats.Rs)
	}
	if got.Stats.Vunit != "C" {
		t.Errorf("GetStats returned wrong vunit: %v", got.Stats.Vunit)
	}
	t.Log(got)
}
