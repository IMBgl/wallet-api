package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"testing"
)

//func TestWalletRepository_Save(t *testing.T) {
//	conn, err := connectDB()
//	if err != nil {
//		t.Errorf("Could not connect to DB %f", err)
//	}
//	repo := New(conn)
//
//	user := domain.NewUser("testName", "testEmail", "testPass")
//
//	err = repo.User().Save(context.Background(), user)
//	if err != nil {
//		t.Errorf("could not save user %f", err)
//	}
//
//	user, err = repo.User().GetByEmail(context.Background(), "testEmail")
//	if err != nil {
//		t.Errorf("could not get user %f", err)
//	}
//
//	t.Logf("User %+v", user)
//}
//
//func TestWalletRepository_SaveUserWithToken(t *testing.T) {
//	conn, err := connectDB()
//	if err != nil {
//		t.Errorf("Could not connect to DB %f", err)
//	}
//	repo := New(conn)
//
//	user := domain.NewUser("testName", "testEmail", "testPass")
//	token := service.NewUserToken("hash", user.Id)
//
//	err = repo.User().SaveUserWithToken(context.Background(), user, token)
//	if err != nil {
//		t.Errorf("could not save user %f", err)
//	}
//
//	user, err = repo.User().GetByEmail(context.Background(), "testEmail")
//	if err != nil {
//		t.Errorf("could not get user %f", err)
//	}
//
//	t.Logf("User %+v", user)
//}

type Person struct {
	Name string
	Sex  bool
}

type Driver struct {
	*Person
}

func Test_mock(t *testing.T) {
	pr := &Person{
		Name: "name",
		Sex:  true,
	}

	dr := Driver{
		Person: pr,
	}

	t.Logf("testing %+v", dr.Person)
}

func connectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, nil
}
