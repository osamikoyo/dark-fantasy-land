package service

import (
	"errors"
	"sync"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"github.com/osamikoyo/dark-fantasy-land/internal/repository"
)

const Limit = 20

func (s *Service) SendArticleToCensor(article *entity.Article) error {
	if article == nil{
		return ErrInvalidInput
	}

	return s.sender.SendToCensor("articles", article)
}

func (s *Service) CreateArticle(article *entity.Article) error {
	if article == nil {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	if err := s.repo.CreateArticle(ctx, article); err != nil {
		return ErrRepositoryFailed
	}

	if err := s.casher.AddArticleToCash(ctx, article); err != nil {
		return ErrCacheSetFailed
	}

	return nil
}

func (s *Service) UpdateArticle(author, title string, update map[string]interface{}) error {
	if author == "" || title == "" {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["author"] = author
	filter["title"] = title

	if err := s.repo.UpdateArticle(ctx, filter, update); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	for key, value := range update {
		if err := s.casher.UpdateArticleInCash(ctx, author, title, key, value); err != nil {
			return ErrCacheSetFailed
		}
	}

	return nil
}

func (s *Service) DeleteArticle(author, title string) error {
	if author == "" || title == "" {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["author"] = author
	filter["title"] = title

	if err := s.repo.DeleteArticle(ctx, filter); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	if err := s.casher.DeleteArticleFromCash(ctx, author, title); err != nil {
		return ErrCacheDelFailed
	}

	return nil
}

func (s *Service) GetOneArticle(author, title string) (*entity.Article, error) {
	if author == "" || title == "" {
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	filter := make(map[string]interface{})
	filter["author"] = author
	filter["title"] = title

	var (
		articleChan = make(chan *entity.Article, 1)
		errChan     = make(chan error, 2)
		wg          sync.WaitGroup
	)

	wg.Add(1)

	go func() {
		defer wg.Done()

		article, err := s.repo.GetArticle(ctx, filter)
		if err != nil {
			errChan <- err
		}
		articleChan <- article
	}()

	go func() {
		defer wg.Done()

		article, err := s.casher.GetArticleFromCash(ctx, author, title)
		if err != nil {
			errChan <- err
		}

		articleChan <- article
	}()

	errCount := 0

	select {
	case article := <-articleChan:
		return article, nil
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

func (s *Service) GetMoreArticles(filter map[string]interface{}) ([]entity.Article, error) {
	ctx, cancel := s.context()
	defer cancel()

	articles, err := s.repo.GetArticlesLimited(ctx, filter, Limit)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrRepositoryFailed
	}

	return articles, nil
}
