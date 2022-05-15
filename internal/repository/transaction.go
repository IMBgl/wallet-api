package repository

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"gorm.io/gorm"
	"time"
)

type TransactionModel struct {
	gorm.Model
	Id         uuid.UUID
	Amount     float32
	Currency   CurrencyValue
	Comment    string
	Type       TransactionTypeValue
	WalletId   uuid.UUID
	UserId     uuid.UUID
	CategoryId uuid.UUID
}

func (TransactionModel) TableName() string {
	return "transactions"
}

func (m *TransactionModel) Entity() *domain.Transaction {
	return domain.NewTransaction(m.Comment, m.Amount, m.Currency.Currency, m.Type.TransactionType, m.UserId, m.CategoryId, m.WalletId)
}

func (m *TransactionModel) FromEntity(e *domain.Transaction) *TransactionModel {
	m.Id = e.Id
	m.Comment = e.Comment
	m.UserId = e.UserId
	m.Currency = CurrencyValue{e.Currency}
	m.Type = TransactionTypeValue{e.Type}
	m.Amount = e.Amount
	m.CategoryId = e.CategoryId
	m.WalletId = e.WalletId
	m.CreatedAt = e.CreatedAt

	return m
}

type transactionRepository struct {
	repository
}

func TransactionRepository(conn *pgx.Conn) *transactionRepository {
	return &transactionRepository{repository{Conn: conn}}
}

func (r *transactionRepository) Save(ctx context.Context, t *domain.Transaction) error {
	_, err := r.Conn.Exec(ctx, "insert into transactions (id, amount, user_id, wallet_id, category_id, currency, \"comment\", \"type\", created_at, updated_at) values($1,$2,$3,$4,$5,$6, $7)",
		t.Id, t.Amount, t.UserId, t.WalletId, t.CategoryId, t.Currency.Val(), t.Comment, t.Type.Val(), t.CreatedAt, time.Now())

	return err
}

func (r *transactionRepository) FindByWalletIdAndUserId(ctx context.Context, walletId, userId uuid.UUID) (list []*domain.Transaction, err error) {
	rows, _ := r.Conn.Query(ctx, "select id, amount, user_id, wallet_id, category_id, currency, \"comment\", \"type\", created_at from transactions where id =$1 and user_id = $2", walletId, userId)

	for rows.Next() {
		i := domain.Transaction{}
		currencyVal := ""
		typeVal := ""

		err := rows.Scan(&i.Id, &i.Amount, &i.UserId, &i.WalletId, &i.CategoryId, &currencyVal, &i.Comment, &typeVal, &i.CreatedAt)
		if err != nil {
			return nil, err
		}

		currency, err := domain.CurrencyFromString(currencyVal)
		if err != nil {
			return nil, err
		}

		transactionType, err := domain.TransactionTypeFromString(typeVal)
		if err != nil {
			return nil, err
		}

		i.Currency = currency
		i.Type = transactionType

		list = append(list, &i)
	}

	return
}

func (r *transactionRepository) SaveAndUpdateWalletBalance(ctx context.Context, t *domain.Transaction, balance float32) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert into transactions (id, amount, user_id, wallet_id, category_id, currency, \"comment\", \"type\", created_at, updated_at,) values($1,$2,$3,$4,$5,$6, $7)",
		t.Id, t.Amount, t.UserId, t.WalletId, t.CategoryId, t.Currency.Val(), t.Comment, t.Type.Val(), t.CreatedAt, time.Now())
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, "update wallets set balance = $1 where id=$2", balance, t.WalletId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
