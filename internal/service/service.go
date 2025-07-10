package service

import (
	"context"
	"time"
)

var (
	
)

type (
	Service struct {
		repo   Repository
		casher Casher

		timeout time.Duration
	}
)

func NewService(repo Repository, casher Casher, timeout time.Duration) *Service {
	return &Service{
		repo: repo,
		casher: casher,
		timeout: timeout,
	}
}

func (s *Service) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.timeout)
}

