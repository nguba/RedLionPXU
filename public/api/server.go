package api

import (
	"context"
	"github.com/nguba/RedLionPXU/internal/device"
	v2 "github.com/nguba/RedLionPXU/public/api/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	v2.UnimplementedRedLionPxuServer
	pid        *device.Pxu
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewServer(pxu *device.Pxu, listener net.Listener) (*Server, error) {
	srv := grpc.NewServer()
	svc := &Server{pid: pxu, listener: listener, grpcServer: srv}
	v2.RegisterRedLionPxuServer(srv, svc)
	return svc, nil
}

func (s *Server) GetStats(_ context.Context, in *v2.GetStatsRequest) (*v2.GetStatsResponse, error) {
	stats, err := s.pid.ReadStats()
	if err != nil {
		return nil, err
	}
	return makeGetStatsResponse(stats), nil
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
	_ = s.listener.Close()
	log.Printf("Stopped gRPC server on: %v", s.listener.Addr())
}

func (s *Server) Start() error {
	go func() {
		log.Printf("Started gRPC server on: %v", s.listener.Addr())
		if err := s.grpcServer.Serve(s.listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return nil
}

func makeGetStatsResponse(stats *device.Stats) *v2.GetStatsResponse {
	return &v2.GetStatsResponse{Stats: &v2.Stats{
		Pv:     stats.Pv,
		Sp:     stats.Sp,
		Out1:   stats.Out1,
		Out2:   stats.Out2,
		At:     stats.At,
		Tp:     stats.TP,
		Ti:     uint32(stats.TI),
		Td:     uint32(stats.TD),
		TGroup: uint32(stats.TGroup),
		Rs:     stats.RS.String(),
		Vunit:  stats.VUnit,
		Pc:     uint32(stats.PC),
		Ps:     uint32(stats.PS),
		Psr:    stats.PSR,
	}}
}
