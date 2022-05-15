package service

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
	"log"
)

type categoryService struct {
	repo Repository
}

type CategoryRepository interface {
	Save(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, c *domain.Category) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Category, error)
	GetChildren(ctx context.Context, c *domain.Category) ([]*domain.Category, error)
	FindByUserIdWithNullParent(ctx context.Context, userId uuid.UUID) ([]*domain.Category, error)
	FindByIdAndUserId(ctx context.Context, id, userId uuid.UUID) (*domain.Category, error)
}

func NewCategoryService(r Repository) *categoryService {
	return &categoryService{repo: r}
}

type CategoryTreeNode struct {
	*domain.Category
	Parent   *domain.Category
	Children []*CategoryTreeNode
}

type CategoryCreateRequest struct {
	Name     string
	Currency domain.Currency
	UserId   uuid.UUID
	ParentId *uuid.UUID
}

type CategoryUpdateRequest struct {
	Name       string
	UserId     uuid.UUID
	CategoryId uuid.UUID
}

type CategoryGetListRequest struct {
	UserId uuid.UUID
}

type CategoryGetOneRequest struct {
	UserId     uuid.UUID
	CategoryId uuid.UUID
}

type CategoryDeleteRequest struct {
	UserId     uuid.UUID
	CategoryId uuid.UUID
}

func (s *categoryService) Create(ctx context.Context, request *CategoryCreateRequest) (category *domain.Category, err error) {
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

	return
}

func (s *categoryService) GetList(ctx context.Context, request *CategoryGetListRequest) (categoryList []*domain.Category, err error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return
	}

	if user == nil {
		err = ErrUserNotFound
		return
	}

	categoryList, err = s.repo.Category().FindByUserIdWithNullParent(ctx, user.Id)
	if err != nil {
		return
	}

	return
}

func (s *categoryService) Delete(ctx context.Context, request *CategoryDeleteRequest) (err error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return
	}

	if user == nil {
		err = ErrUserNotFound
		return
	}

	category, err := s.repo.Category().FindByIdAndUserId(ctx, request.CategoryId, user.Id)
	if err != nil {
		return
	}

	if category == nil {
		return ErrCategoryNotFound
	}

	err = s.repo.Category().Delete(ctx, category)
	if err != nil {
		return
	}

	return
}

func (s *categoryService) Update(ctx context.Context, request *CategoryUpdateRequest) (node *CategoryTreeNode, err error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	category, err := s.repo.Category().FindByIdAndUserId(ctx, request.CategoryId, user.Id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, ErrCategoryNotFound
	}

	category.Name = request.Name

	err = s.repo.Category().Save(ctx, category)
	if err != nil {
		return nil, err
	}

	return s.getCategoryTree(ctx, category)
}

func (s *categoryService) GetOne(ctx context.Context, request *CategoryGetOneRequest) (*CategoryTreeNode, error) {
	node := &CategoryTreeNode{}

	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		err = ErrUserNotFound
		return nil, err
	}

	category, err := s.repo.Category().FindByIdAndUserId(ctx, request.CategoryId, user.Id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, ErrCategoryNotFound
	}

	node.Category = category

	log.Printf("get category param %+v", category)
	log.Println()
	node, err = s.getCategoryTree(ctx, category)
	log.Printf("category tree %+v", category)
	log.Println()
	return node, nil
}

func (s *categoryService) getCategoryTree(ctx context.Context, c *domain.Category) (*CategoryTreeNode, error) {
	node := &CategoryTreeNode{}
	node.Category = c

	if node.Category.ParentId != nil {
		parent, err := s.repo.Category().FindByIdAndUserId(ctx, *node.ParentId, node.UserId)
		if err != nil {
			return nil, err
		}
		node.Parent = parent
	}

	children, err := s.repo.Category().GetChildren(ctx, node.Category)
	log.Printf("children %+v", children)
	log.Println()

	if err != nil {
		return nil, err
	}

	if len(children) > 0 {
		for _, child := range children {
			childNode, err := s.getCategoryTree(ctx, child)
			if err != nil {
				return nil, err
			}

			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}
