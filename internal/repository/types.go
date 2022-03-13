package repository

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/IMBgl/go-wallet-api/internal/domain"
)

type CurrencyValue struct {
	Currency domain.Currency
}

func (t *CurrencyValue) Scan(value interface{}) error {
	stringVal, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to decode CurrencyValue:", value))
	}

	val, err := domain.CurrencyFromString(stringVal)
	t.Currency = val
	return err
}

func (t CurrencyValue) Value() (driver.Value, error) {
	return t.Currency.Val(), nil
}

type TransactionTypeValue struct {
	TransactionType domain.TransactionType
}

func (t *TransactionTypeValue) Scan(value interface{}) error {
	stringVal, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to decode CurrencyValue:", value))
	}

	val, err := domain.TransactionTypeFromString(stringVal)
	t.TransactionType = val
	return err
}

func (t TransactionTypeValue) Value() (driver.Value, error) {
	return t.TransactionType.Val(), nil
}
