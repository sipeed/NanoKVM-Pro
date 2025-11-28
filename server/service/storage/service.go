package storage

import "os"

type Service struct{}

func NewService() *Service {
	_ = os.Remove(sentinelPath)
	return &Service{}
}
