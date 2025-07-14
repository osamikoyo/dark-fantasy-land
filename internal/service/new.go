package service

import (
	"errors"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"github.com/osamikoyo/dark-fantasy-land/internal/repository"
)

func (s *Service) CreateNew(new *entity.New) error {
	if new == nil{
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	if err := s.repo.CreateNew(ctx, new);err != nil{
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrAlreadyExists
		}

		return ErrRepositoryFailed
	}

	if err := s.casher.AddNewToCash(ctx, new);err != nil{
		return ErrCacheSetFailed
	}

	return nil
}

func (s *Service) UpdateNew(author, title string, update map[string]interface{}) error {
	if author == "" || title == "" || update == nil {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["author"] = author
	filter["title"] = title

	if err := s.repo.UpdateNew(ctx, filter, update);err != nil{
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	for key, value := range update{
		if err := s.casher.UpdateNewInCash(ctx, author, title, key, value);err != nil{
			return ErrCacheSetFailed
		}
	}

	return nil
}