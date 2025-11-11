package article

import (
	"context"

	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/response"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateArticleRequest) (Article, error) {
	// create article
	return s.repo.Insert(ctx, req.Title, req.Content, req.Category, req.Status)
}

func (s *Service) List(ctx context.Context, limit, page int, filter ListFilter) ([]Article, response.Meta, error) {
	// set default limit and page
	newLimit := max(limit, 10)
	newPage := max(page, 1)

	// calculate offset
	offset := (newPage - 1) * newLimit
	items, total, err := s.repo.List(ctx, newLimit, offset, filter)
	if err != nil {
		return []Article{}, response.Meta{}, err
	}

	return items, *response.PageMeta(newLimit, offset, total), nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (Article, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateArticleRequest) (Article, error) {
	// get current article
	curr, err := s.repo.FindByID(ctx, id)
	// check if article exists
	if err != nil {
		return Article{}, err
	}

	if req.Title != "" {
		curr.Title = req.Title
	}
	if req.Content != "" {
		curr.Content = req.Content
	}
	if req.Category != "" {
		curr.Category = req.Category
	}
	if req.Status != "" {
		curr.Status = req.Status
	}

	// update article
	up := UpdateArticleRequest{
		Title:    curr.Title,
		Content:  curr.Content,
		Category: curr.Category,
		Status:   curr.Status,
	}

	return s.repo.UpdateAll(ctx, id, up.Title, up.Content, up.Category, up.Status)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	// delete article
	return s.repo.Delete(ctx, id)
}
