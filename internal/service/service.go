package service

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
)

type service struct {
	repo        Repository
	user        UserService
	token       TokenService
	wallet      WalletService
	category    CategoryService
	transaction TransactionService
}

type Service interface {
	User() UserService
	Token() TokenService
	Wallet() WalletService
	Category() CategoryService
	Transaction() TransactionService
}

type Repository interface {
	User() UserRepository
	Token() TokenRepository
	Wallet() WalletRepository
	Category() CategoryRepository
	Transaction() TransactionRepository
}

type UserService interface {
	SingUp(ctx context.Context, request SignUpRequest) (*domain.User, *UserToken, error)
	SingIn(ctx context.Context, request SignInRequest) (*domain.User, *UserToken, error)
}

type TokenService interface {
	CreateForUser(u *domain.User) *UserToken
	GetByValue(ctx context.Context, v string) (*UserToken, error)
}

type WalletService interface {
	Create(ctx context.Context, request *WalletCreateRequest) (*domain.Wallet, error)
	GetList(ctx context.Context, request *WalletGetListRequest) (walletList []*domain.Wallet, err error)
	Delete(ctx context.Context, request *WalletDeleteRequest) error
	Update(ctx context.Context, request *WalletUpdateRequest) (*domain.Wallet, error)
}

type CategoryService interface {
	Create(ctx context.Context, request *CategoryCreateRequest) (category *domain.Category, parent *domain.Category, err error)
}

type TransactionService interface {
	Create(ctx context.Context, request *TransactionCreateRequest) (*domain.Transaction, error)
}

func (s *service) User() UserService {
	return s.user
}

func (s *service) Token() TokenService {
	return s.token
}

func (s *service) Wallet() WalletService {
	return s.wallet
}

func (s *service) Category() CategoryService {
	return s.category
}

func (s *service) Transaction() TransactionService {
	return s.transaction
}

func New(repo Repository) *service {
	ts := &tokenServiсe{repo: repo}
	us := NewUserService(repo, ts)
	ws := NewWalletService(repo)
	cs := NewCategoryService(repo)
	trs := NewTransactionService(repo)

	return &service{
		repo:        repo,
		token:       ts,
		user:        us,
		wallet:      ws,
		category:    cs,
		transaction: trs,
	}
}
