package service

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
)

type categoryService struct {
	repo Repository
}

type CategoryRepository interface {
	Save(ctx context.Context, c *domain.Category) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.Category, error)
}

func NewCategoryService(r Repository) *categoryService {
	return &categoryService{repo: r}
}

type CategoryCreateRequest struct {
	Name     string
	Currency domain.Currency
	UserId   uuid.UUID
	ParentId *uuid.UUID
}

func (s *categoryService) Create(ctx context.Context, request *CategoryCreateRequest) (category *domain.Category, parent *domain.Category, err error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return
	}

	if user == nil {
		err = ErrUserNotFound
		return
	}

	category = domain.NewCategory(request.Name, request.Currency, request.UserId)
	category.ParentId = request.ParentId

	err = s.repo.Category().Save(ctx, category)
	if err != nil {
		return
	}

	parent, err = s.repo.Category().GetById(ctx, *category.ParentId)
	if err != nil {
		return
	}

	return
}
