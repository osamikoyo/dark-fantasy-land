package service

import (
	"context"
	"errors"
	"time"
)

const (
	Limit = 20
	RetrierAttemps = 3
	RetrierDuration = 2 * time.Second
)

var (
	ErrNotFound         = errors.New("not found")
	ErrTimeout          = errors.New("operation timeout")
	ErrInvalidInput     = errors.New("invalid input")
	ErrAlreadyExists    = errors.New("already exist")
	ErrCacheSetFailed   = errors.New("cache set failed")
	ErrCacheGetFailed   = errors.New("cache get failed")
	ErrCacheDelFailed   = errors.New("cache delete failed")
	ErrRepositoryFailed = errors.New("repository operation failed")
	ErrInternal         = errors.New("internal service error")
)

type (
	Service struct {
		repo   Repository
		casher Casher
		sender Sender

		timeout time.Duration
	}
)

func NewService(repo Repository, casher Casher, sender Sender, timeout time.Duration) *Service {
	return &Service{
		repo:    repo,
		casher:  casher,
		timeout: timeout,
	}
}

func (s *Service) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.timeout)
}

func (s *Service) sendToCensor(value interface{}, subj string) error {
	if value == nil {
		return ErrInvalidInput
	}

	if err := s.sender.SendToCensor(subj, value); err != nil {
		return ErrInternal
	}

	return nil
}
