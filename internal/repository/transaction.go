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
	_, err := r.Conn.Exec(ctx, "insert into transactions (id, amount, user_id, wallet_id, category_id, currency, \"comment\", \"type\", created_at, updated_at,) values($1,$2,$3,$4,$5,$6, $7)",
		t.Id, t.Amount, t.UserId, t.WalletId, t.CategoryId, t.Currency.Val(), t.Comment, t.Type.Val(), t.CreatedAt, time.Now())

	return err
}
