package pxu

import (
	"context"
	"github.com/nguba/RedLionPXU/internal/device"
	v1 "github.com/nguba/RedLionPXU/public/api/pxu/v1"
	"time"
)

type PxuServer struct {
	v1.UnimplementedRedLionPxuServer
	pid *device.Pxu
}

func NewPxuServer(unitId device.UnitId, modbus device.Modbus) (*PxuServer, error) {
	pxu, err := device.NewPxu(unitId, modbus, time.Second, 3)
	if err != nil {
		return nil, err
	}
	return &PxuServer{pid: pxu}, nil
}

func (s *PxuServer) GetStats(_ context.Context, in *v1.GetStatsRequest) (*v1.GetStatsResponse, error) {
	stats, err := s.pid.ReadStats()
	if err != nil {
		return nil, err
	}
	return makeGetStatsResponse(stats), nil
}

func makeGetStatsResponse(stats *device.Stats) *v1.GetStatsResponse {
	return &v1.GetStatsResponse{Stats: &v1.Stats{
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
