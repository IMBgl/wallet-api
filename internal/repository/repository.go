package repository

import (
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/jackc/pgx/v4"
)

type repository struct {
	Conn        *pgx.Conn
	user        *userRepository
	token       *tokenRepository
	wallet      *walletRepository
	category    *categoryRepository
	transaction *transactionRepository
}

//func (r *repository) AsTransaction(ctx context.Context, payload func() error) error {
//	r.DB = r.DB.WithContext(ctx).Begin(&sql.TxOptions{})
//	defer func() {
//		r.DB.Rollback()
//	}()
//	err := payload()
//
//	if err != nil {
//		log.Println("rolling back")
//		r.DB.Rollback()
//		return err
//	}
//
//	return r.DB.Commit().Error
//}

func (r *repository) User() service.UserRepository {
	return r.user
}

func (r *repository) Token() service.TokenRepository {
	return r.token
}

func (r *repository) Wallet() service.WalletRepository {
	return r.wallet
}

func (r *repository) Category() service.CategoryRepository {
	return r.category
}

func (r *repository) Transaction() service.TransactionRepository {
	return r.transaction
}

func New(conn *pgx.Conn) *repository {
	return &repository{
		Conn:        conn,
		user:        UserRepository(conn),
		token:       TokenRepository(conn),
		wallet:      WalletRepository(conn),
		category:    CategoryRepository(conn),
		transaction: TransactionRepository(conn),
	}
}
