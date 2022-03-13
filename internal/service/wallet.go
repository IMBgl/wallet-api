package service

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
)

type walletService struct {
	repo Repository
}

type WalletCreateRequest struct {
	Name     string
	UserId   uuid.UUID
	Currency string
	Balance  float32
}

type WalletUpdateRequest struct {
	Name     string
	UserId   uuid.UUID
	WalletId uuid.UUID
}

type WalletGetListRequest struct {
	UserId uuid.UUID
}

type WalletDeleteRequest struct {
	WalletId uuid.UUID
	UserId   uuid.UUID
}

type WalletRepository interface {
	Save(ctx context.Context, w *domain.Wallet) error
	Delete(ctx context.Context, w *domain.Wallet) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	GetByIdAndUserId(ctx context.Context, id, userId uuid.UUID) (*domain.Wallet, error)
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Wallet, error)
}

func NewWalletService(r Repository) *walletService {
	return &walletService{repo: r}
}

func (s *walletService) Create(ctx context.Context, request *WalletCreateRequest) (*domain.Wallet, error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	currency, err := domain.CurrencyFromString(request.Currency)
	if err != nil {
		return nil, err
	}

	wallet := &domain.Wallet{
		Id:       uuid.New(),
		Name:     request.Name,
		Currency: currency,
		Balance:  request.Balance,
		UserId:   request.UserId,
	}

	err = s.repo.Wallet().Save(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) GetList(ctx context.Context, request *WalletGetListRequest) (walletList []*domain.Wallet, err error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return
	}

	if user == nil {
		err = ErrUserNotFound
		return
	}

	walletList, err = s.repo.Wallet().FindByUserId(ctx, user.Id)

	return
}

func (s *walletService) Delete(ctx context.Context, request *WalletDeleteRequest) error {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return err
	}

	if user == nil {
		err = ErrUserNotFound
		return err
	}

	wallet, err := s.repo.Wallet().GetByIdAndUserId(ctx, request.WalletId, user.Id)
	if err != nil {
		return err
	}

	if wallet == nil {
		return ErrWalletNotFound
	}

	err = s.repo.Wallet().Delete(ctx, wallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *walletService) Update(ctx context.Context, request *WalletUpdateRequest) (*domain.Wallet, error) {
	user, err := s.repo.User().GetById(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	wallet, err := s.repo.Wallet().GetByIdAndUserId(ctx, request.WalletId, request.UserId)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	wallet.Name = request.Name

	err = s.repo.Wallet().Save(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
