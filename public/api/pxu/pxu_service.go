package pxu

import (
	"context"
	v1 "github.com/nguba/RedLionPXU/public/api/pxu/v1"
	"log"
	"sync"
)

type PxuService struct {
	v1.UnimplementedRedLionPxuServiceServer
	mu sync.Mutex
}

func NewPxuService() *PxuService {
	return &PxuService{}
}

func (s *PxuService) GetStats(_ context.Context, in *v1.GetStatsRequest) (*v1.GetStatsResponse, error) {
	log.Printf("RCV: %v", in.String())
	return nil, nil
}
