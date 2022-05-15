package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Budget struct {
	Id         uuid.UUID
	Name       string
	CategoryId uuid.UUID
	Limit      float32
	Amount     float32
}

type Category struct {
	Id        uuid.UUID
	Name      string
	UserId    uuid.UUID
	Currency  Currency
	CreatedAt time.Time
	ParentId  *uuid.UUID
}

type Transaction struct {
	Id         uuid.UUID
	Amount     float32
	Currency   Currency
	Comment    string
	Type       TransactionType
	WalletId   uuid.UUID
	UserId     uuid.UUID
	CategoryId uuid.UUID
	CreatedAt  time.Time
}

type Transfer struct {
	Id               uuid.UUID
	OutTransactionId uuid.UUID
	InTransactionId  uuid.UUID
	CreatedAt        time.Time
}

type Wallet struct {
	Id        uuid.UUID
	Name      string
	UserId    uuid.UUID
	Balance   float32
	Currency  Currency
	CreatedAt time.Time
}

func (w *Wallet) updateBalance(t *Transaction) error {
	if t.Type.IsIn() {
		w.Balance += t.Amount
	} else if t.Type.IsOut() {
		w.Balance -= t.Amount
	} else {
		return ErrInvalidTransactionType
	}

	return nil
}

func NewUser(name, email, password string) *User {
	return &User{
		Id:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}

func NewCategory(name string, currency Currency, userId uuid.UUID) *Category {
	return &Category{
		Id:        uuid.New(),
		Name:      name,
		Currency:  currency,
		UserId:    userId,
		CreatedAt: time.Now(),
		ParentId:  nil,
	}
}

func NewTransaction(comment string, amount float32, currency Currency, transactionType TransactionType, userId, categoryId, walletId uuid.UUID) *Transaction {
	return &Transaction{
		Id:         uuid.New(),
		Amount:     amount,
		Comment:    comment,
		Currency:   currency,
		Type:       transactionType,
		UserId:     userId,
		CategoryId: categoryId,
		WalletId:   walletId,
		CreatedAt:  time.Now(),
	}
}
