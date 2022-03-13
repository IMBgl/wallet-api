package domain

type DomainError struct {
	message string
}

func (d *DomainError) Error() string {
	return d.message
}

func NewError(message string) error {
	return &DomainError{message: message}
}

var (
	ErrWalletNotFound      = NewError("Wallet not found")
	ErrCategoryNotFound    = NewError("Category not found")
	ErrUserNotFound        = NewError("User not found")
	ErrTransactionNotFound = NewError("Transaction not found")

	ErrTransactionWalletCurrencyMismatch   = NewError("Transaction and Wallet must have same currency")
	ErrTransactionCategoryCurrencyMismatch = NewError("Transaction and Category must have same currency")

	ErrInvalidCurrency = NewError("Invalid currency")

	ErrInvalidTransactionType = NewError("Invalid transaction type")
)
