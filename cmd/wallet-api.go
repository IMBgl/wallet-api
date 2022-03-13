package main

import (
	"context"
	"fmt"
	"github.com/IMBgl/go-wallet-api/internal/handler"
	"github.com/IMBgl/go-wallet-api/internal/repository"
	"github.com/IMBgl/go-wallet-api/internal/service"
	pgx "github.com/jackc/pgx/v4"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	conn, err := connectDB()
	if err != nil {
		log.Printf("Could not connect to database %v", err)
	}

	repo := repository.New(conn)
	srv := service.New(repo)

	router := handler.ApiHandler(srv).Routes()

	err = http.ListenAndServe(os.Getenv("APP_HOST"), router)
	if err != nil {
		fmt.Printf("Server error %v", err)
	}
}

func connectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, nil
}
