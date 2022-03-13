package service

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
)

type transactionService struct {
	repo Repository
}

type TransactionCreateRequest struct {
	UserId          uuid.UUID
	WalletId        uuid.UUID
	CategoryId      uuid.UUID
	Currency        domain.Currency
	Comment         string
	Amount          float32
	TransactionType domain.TransactionType
}

type TransactionRepository interface {
	Save(ctx context.Context, w *domain.Transaction) error
}

func NewTransactionService(r Repository) TransactionService {
	return &transactionService{repo: r}
}

func (s *transactionService) Create(ctx context.Context, request *TransactionCreateRequest) (*domain.Transaction, error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	category, err := s.repo.Category().GetById(ctx, request.CategoryId)
	if err != nil {
		return nil, err
	}

	wallet, err := s.repo.Wallet().GetById(ctx, request.WalletId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}
	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}
	if wallet == nil {
		return nil, domain.ErrWalletNotFound
	}

	transaction := domain.NewTransaction(request.Comment, request.Amount, request.Currency, request.TransactionType, request.UserId, request.CategoryId, request.WalletId)

	err = s.repo.Transaction().Save(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
