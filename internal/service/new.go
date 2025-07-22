package service

import (
	"errors"
	"sync"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"github.com/osamikoyo/dark-fantasy-land/internal/repository"
)

func (s *Service) CreateNew(new *entity.New) error {
	if new == nil {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	if err := s.repo.CreateNew(ctx, new); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrAlreadyExists
		}

		return ErrRepositoryFailed
	}

	if err := s.sendToCensor(new, "news"); err != nil {
		return err
	}

	if err := s.casher.AddNewToCash(ctx, new); err != nil {
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

	if err := s.repo.UpdateNew(ctx, filter, update); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	for key, value := range update {
		if err := s.casher.UpdateNewInCash(ctx, author, title, key, value); err != nil {
			return ErrCacheSetFailed
		}
	}

	return nil
}

func (s *Service) DeleteNew(author, title string) error {
	if author == "" || title == "" {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["author"] = author
	filter["title"] = title

	if err := s.repo.DeleteNew(ctx, filter); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	if err := s.casher.DeleteNewFromCash(ctx, author, title); err != nil {
		return ErrCacheDelFailed
	}

	return nil
}

func (s *Service) GetOneNew(author, title string) (*entity.New, error) {
	if author == "" || title == "" {
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["author"] = author
	filter["title"] = title

	var (
		wg      sync.WaitGroup
		newChan chan *entity.New
		errChan chan error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		new, err := s.repo.GetNew(ctx, filter)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				err = ErrNotFound
			}

			errChan <- err

			return
		}

		newChan <- new
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		new, err := s.casher.GetNewFromCash(ctx, author, title)
		if err != nil {
			errChan <- ErrCacheGetFailed

			return
		}

		newChan <- new
	}()

	errCount := 0

	for errCount != 2 {
		select {
		case new := <-newChan:
			return new, nil
		case <-errChan:
			errCount++
		}
	}

	return nil, ErrInternal
}

func (s *Service) GetManyNew(filter map[string]interface{}) ([]entity.New, error) {
	if filter == nil {
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	news, err := s.repo.GetNewsLimited(ctx, filter, Limit)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrRepositoryFailed
	}

	return news, nil
}
