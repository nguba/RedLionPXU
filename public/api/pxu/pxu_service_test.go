package pxu

import (
	"context"
	v1 "github.com/nguba/RedLionPXU/public/api/pxu/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// initialize a listener that doesn't depend on a networked implementation
func init() {
	lis = bufconn.Listen(bufSize)
}

func setupTestServer(t *testing.T, svc v1.RedLionPxuServiceServer) v1.RedLionPxuServiceClient {
	t.Helper() // Mark this as a helper function

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { lis.Close() })

	srv := grpc.NewServer()
	t.Cleanup(func() { srv.Stop() })

	v1.RegisterRedLionPxuServiceServer(srv, svc)

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
	t.Cleanup(func() { conn.Close() })

	return v1.NewRedLionPxuServiceClient(conn)
}

func TestApi_GetStats(t *testing.T) {
	client := setupTestServer(t, NewPxuService())
	stats, err := client.GetStats(context.Background(), &v1.GetStatsRequest{})
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	t.Log(stats)
}
