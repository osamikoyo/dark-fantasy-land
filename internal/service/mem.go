package service

import (
	"errors"
	"sync"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"github.com/osamikoyo/dark-fantasy-land/internal/repository"
)

func (s *Service) SendMemToCensor(mem *entity.Mem) error {
	if mem == nil {
		return ErrInvalidInput
	}

	if err := s.sender.SendToCensor("mems", mem); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *Service) CreateMem(mem *entity.Mem) error {
	if mem == nil {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	if err := s.repo.CreateMem(ctx, mem); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrAlreadyExists
		}

		return ErrRepositoryFailed
	}

	if err := s.casher.AddMemToCash(ctx, mem); err != nil {
		return ErrCacheSetFailed
	}

	return nil
}

func (s *Service) UpdateMem(image_name, author string, update map[string]interface{}) error {
	if author == "" || image_name == "" || update == nil {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["image_name"] = image_name
	filter["author"] = author

	if err := s.repo.UpdateMem(ctx, filter, update); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	for key, value := range update {
		if err := s.casher.UpdateMemInCash(ctx, image_name, author, key, value); err != nil {
			return ErrCacheSetFailed
		}
	}

	return nil
}

func (s *Service) GetOneMem(image_name, author string) (*entity.Mem, error) {
	if author == "" || image_name == "" {
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["image_name"] = image_name
	filter["author"] = author

	var (
		memChan = make(chan *entity.Mem, 1)
		errChan = make(chan error, 2)
		wg      sync.WaitGroup
	)

	wg.Add(1)

	go func() {
		defer wg.Done()

		mem, err := s.repo.GetMem(ctx, filter)
		if err != nil {
			errChan <- err
		}
		memChan <- mem
	}()

	go func() {
		defer wg.Done()

		mem, err := s.casher.GetMemFromCash(ctx, image_name, author)
		if err != nil {
			errChan <- err
		}

		memChan <- mem
	}()

	errCount := 0

	select {
	case mem := <-memChan:
		return mem, nil
	case err := <-errChan:
		errCount++
		if errCount < 1 {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}

			return nil, err
		}
	}

	return nil, ErrInternal
}

func (s *Service) GetManyMems(filter map[string]interface{}) ([]entity.Mem, error) {
	if filter == nil {
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	mems, err := s.repo.GetMemsLimited(ctx, filter, Limit)
	if err != nil{
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrRepositoryFailed
	}

	return mems, nil
}

func (s *Service) DeleteMem(image_name, author string) error {
	if author == "" || image_name == "" {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["image_name"] = image_name
	filter["author"] = author

	if err := s.repo.DeleteArticle(ctx, filter);err != nil{
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	if err := s.casher.DeleteMemInCash(ctx, image_name, author);err != nil{
		return ErrCacheDelFailed
	}

	return nil
}