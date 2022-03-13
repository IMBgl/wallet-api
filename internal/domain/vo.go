package domain

import "strings"

const currencyRUR = "rur"
const currencyEUR = "eur"
const currencyUSD = "usd"

const transactionIn = "in"
const transactionOut = "out"

type Currency struct {
	value string
}

func CurrencyRUR() Currency {
	return Currency{value: currencyRUR}
}

func CurrencyUSD() Currency {
	return Currency{value: currencyUSD}
}

func CurrencyEUR() Currency {
	return Currency{value: currencyEUR}
}

func (c *Currency) Equals(cur *Currency) bool {
	return c.value == cur.value
}

func (c *Currency) Val() string {
	return c.value
}

func CurrencyFromString(val string) (Currency, error) {
	val = strings.ToLower(val)

	if val == currencyRUR {
		return CurrencyRUR(), nil
	} else if val == currencyEUR {
		return CurrencyEUR(), nil
	} else if val == currencyUSD {
		return CurrencyUSD(), nil
	}

	return Currency{}, ErrInvalidCurrency
}

type TransactionType struct {
	value string
}

func TransactionTypeIn() TransactionType {
	return TransactionType{value: transactionIn}
}

func TransactionTypeOut() TransactionType {
	return TransactionType{value: transactionOut}
}

func (tt *TransactionType) Val() string {
	return tt.value
}

func (tt *TransactionType) IsIn() bool {
	return tt.value == transactionIn
}

func (tt *TransactionType) IsOut() bool {
	return tt.value == transactionOut
}

func (tt *TransactionType) Equals(ett *TransactionType) bool {
	return tt.value == ett.value
}

func TransactionTypeFromString(val string) (TransactionType, error) {
	val = strings.ToLower(val)

	if val == transactionIn {
		return TransactionTypeIn(), nil
	} else if val == transactionOut {
		return TransactionTypeOut(), nil
	}

	return TransactionType{}, ErrInvalidTransactionType
}
