package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryDto struct {
	Name     string
	UserId   uuid.UUID
	Total    float32
	Currency Currency
}

type CreateTransactionDto struct {
	Amount    float32
	Currency  Currency
	Comment   string
	Type      TransactionType
	WalletId  uuid.UUID
	UserId    uuid.UUID
	Timestamp time.Time
}

type CreateWalletDto struct {
	Name     string
	UserId   uuid.UUID
	Balance  float32
	Currency Currency
}

type CreateUserDto struct {
	Name     string
	Email    string
	Password string
}
