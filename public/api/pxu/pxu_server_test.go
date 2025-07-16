package pxu

import (
	"context"
	"github.com/nguba/RedLionPXU/internal/device"
	v1 "github.com/nguba/RedLionPXU/public/api/pxu/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
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

func setupTestServer(t *testing.T) v1.RedLionPxuClient {
	t.Helper() // Mark this as a helper function

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { lis.Close() })

	srv := grpc.NewServer()
	t.Cleanup(func() { srv.Stop() })

	svc, err := NewPxuServer(unit, modbus)
	v1.RegisterRedLionPxuServer(srv, svc)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

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

	return v1.NewRedLionPxuClient(conn)
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

	got, err := client.GetStats(context.Background(), &v1.GetStatsRequest{})
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
	if got.Stats.Rs != device.RunStatusRun.String() {
		t.Errorf("GetStats returned wrong rs: %v", got.Stats.Rs)
	}
	if got.Stats.Vunit != "C" {
		t.Errorf("GetStats returned wrong vunit: %v", got.Stats.Vunit)
	}
}
