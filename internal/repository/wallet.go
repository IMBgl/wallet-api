package repository

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"gorm.io/gorm"
	"time"
)

type WalletModel struct {
	gorm.Model
	Id       uuid.UUID
	Name     string
	UserId   uuid.UUID
	Currency CurrencyValue
	Balance  float32
}

func (WalletModel) TableName() string {
	return "wallets"
}

func (m *WalletModel) Entity() (*domain.Wallet, error) {
	return &domain.Wallet{
		m.Id,
		m.Name,
		m.UserId,
		m.Balance,
		m.Currency.Currency,
		m.CreatedAt,
	}, nil
}

func (m *WalletModel) FromEntity(e *domain.Wallet) *WalletModel {
	m.Id = e.Id
	m.Name = e.Name
	m.UserId = e.UserId
	m.Balance = e.Balance
	m.Currency = CurrencyValue{e.Currency}

	return m
}

type walletRepository struct {
	repository
}

func WalletRepository(conn *pgx.Conn) *walletRepository {
	return &walletRepository{repository{Conn: conn}}
}

func (r *walletRepository) Save(ctx context.Context, w *domain.Wallet) error {
	_, err := r.Conn.Exec(ctx, `insert into wallets (id, "name", user_id, currency,balance, created_at, updated_at)
									values($1,$2,$3,$4,$5,$6, $7)
									on conflict (id) do update 
									set name = $2, updated_at = $7;`, w.Id, w.Name, w.UserId, w.Currency.Val(), w.Balance, w.CreatedAt, time.Now())

	return err
}

func (r *walletRepository) Delete(ctx context.Context, w *domain.Wallet) error {
	_, err := r.Conn.Exec(ctx, "delete from wallets where id=$1", w.Id)

	return err
}

func (r *walletRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	wallet := domain.Wallet{}
	currencyVal := ""

	err := r.Conn.QueryRow(ctx, "select id, \"name\", user_id, currency,balance, created_at from wallets where id=$1", id).Scan(&wallet.Id, &wallet.Name, &wallet.UserId, &currencyVal, &wallet.Balance, &wallet.CreatedAt)
	if err != nil {
		return nil, err
	}
	currency, err := domain.CurrencyFromString(currencyVal)
	if err != nil {
		return nil, err
	}

	wallet.Currency = currency

	return &wallet, nil
}

func (r *walletRepository) GetByIdAndUserId(ctx context.Context, id, userId uuid.UUID) (*domain.Wallet, error) {
	wallet := domain.Wallet{}
	currencyVal := ""

	err := r.Conn.QueryRow(ctx, "select id, \"name\", user_id, currency,balance, created_at from wallets where id=$1 and user_id=$2", id, userId).Scan(&wallet.Id, &wallet.Name, &wallet.UserId, &currencyVal, &wallet.Balance, &wallet.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	currency, err := domain.CurrencyFromString(currencyVal)
	if err != nil {
		return nil, err
	}

	wallet.Currency = currency

	return &wallet, nil
}

func (r *walletRepository) FindByUserId(ctx context.Context, userId uuid.UUID) (list []*domain.Wallet, err error) {
	rows, _ := r.Conn.Query(ctx, "select id, \"name\", user_id, currency,balance, created_at from wallets where user_id=$1", userId)

	for rows.Next() {
		i := domain.Wallet{}
		currencyVal := ""

		err := rows.Scan(&i.Id, &i.Name, &i.UserId, &currencyVal, &i.Balance, &i.CreatedAt)
		if err != nil {
			return nil, err
		}

		currency, err := domain.CurrencyFromString(currencyVal)
		if err != nil {
			return nil, err
		}

		i.Currency = currency
		list = append(list, &i)
	}

	return
}
