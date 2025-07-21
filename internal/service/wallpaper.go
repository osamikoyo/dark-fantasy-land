package service

import (
	"errors"
	"sync"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"github.com/osamikoyo/dark-fantasy-land/internal/repository"
)

func (s *Service) CreateWallpaper(wallpaper *entity.Wallpaper) error {
	if wallpaper == nil {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	if err := s.repo.CreateWallpaper(ctx, wallpaper); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrAlreadyExists
		}
		return ErrRepositoryFailed
	}

	if err := s.sendToCensor(wallpaper, "wallpapers");err != nil{
		return err
	}

	if err := s.casher.AddWallpaperToCash(ctx, wallpaper); err != nil {
		return ErrCacheSetFailed
	}

	return nil
}

func (s *Service) UpdateWallpaper(imageName, title string, update map[string]interface{}) error {
	if imageName == "" || title == "" || update == nil {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := map[string]interface{}{
		"image_name": imageName,
		"title":      title,
	}

	if err := s.repo.UpdateWallpaper(ctx, filter, update); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return ErrRepositoryFailed
	}

	for key, value := range update {
		if err := s.casher.UpdateWallpaperInCash(ctx, imageName, title, key, value); err != nil {
			return ErrCacheSetFailed
		}
	}

	return nil
}

func (s *Service) DeleteWallpaper(imageName, title string) error {
	if imageName == "" || title == "" {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := map[string]interface{}{
		"image_name": imageName,
		"title":      title,
	}

	if err := s.repo.DeleteWallpaper(ctx, filter); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return ErrRepositoryFailed
	}

	if err := s.casher.DeleteWallpaperFromCash(ctx, imageName, title); err != nil {
		return ErrCacheDelFailed
	}

	return nil
}

func (s *Service) GetOneWallpaper(imageName, title string) (*entity.Wallpaper, error) {
	if imageName == "" || title == "" {
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := map[string]interface{}{
		"image_name": imageName,
		"title":      title,
	}

	var (
		wg            sync.WaitGroup
		wallpaperChan = make(chan *entity.Wallpaper, 1)
		errChan       = make(chan error, 2)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		wallpaper, err := s.repo.GetWallpaper(ctx, filter)
		if err != nil {
			errChan <- err
			return
		}
		wallpaperChan <- wallpaper
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		wallpaper, err := s.casher.GetWallpaperFromCash(ctx, imageName, title)
		if err != nil {
			errChan <- ErrCacheGetFailed
			return
		}
		wallpaperChan <- wallpaper
	}()

	errCount := 0
	for errCount != 2 {
		select {
		case wallpaper := <-wallpaperChan:
			return wallpaper, nil
		case <-errChan:
			errCount++
		}
	}

	return nil, ErrInternal
}

func (s *Service) GetManyWallpapers(filter map[string]interface{}) ([]entity.Wallpaper, error) {
	if filter == nil {
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	wallpapers, err := s.repo.GetWallpapersLimited(ctx, filter, Limit)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrRepositoryFailed
	}

	return wallpapers, nil
}
