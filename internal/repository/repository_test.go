package repository

import (
	"context"
	"fmt"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/jackc/pgx/v4"
	"os"
	"testing"
)

func TestWalletRepository_Save(t *testing.T) {
	conn, err := connectDB()
	if err != nil {
		t.Errorf("Could not connect to DB %f", err)
	}
	repo := New(conn)

	user := domain.NewUser("testName", "testEmail", "testPass")

	err = repo.User().Save(context.Background(), user)
	if err != nil {
		t.Errorf("could not save user %f", err)
	}

	user, err = repo.User().GetByEmail(context.Background(), "testEmail")
	if err != nil {
		t.Errorf("could not get user %f", err)
	}

	t.Logf("User %+v", user)
}

func TestWalletRepository_SaveUserWithToken(t *testing.T) {
	conn, err := connectDB()
	if err != nil {
		t.Errorf("Could not connect to DB %f", err)
	}
	repo := New(conn)

	user := domain.NewUser("testName", "testEmail", "testPass")
	token := service.NewUserToken("hash", user.Id)

	err = repo.User().SaveUserWithToken(context.Background(), user, token)
	if err != nil {
		t.Errorf("could not save user %f", err)
	}

	user, err = repo.User().GetByEmail(context.Background(), "testEmail")
	if err != nil {
		t.Errorf("could not get user %f", err)
	}

	t.Logf("User %+v", user)
}

func connectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, nil
}
