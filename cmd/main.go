package main

import (
	"banking_system_golang/internal/handlers"
	"banking_system_golang/internal/repository"
	"banking_system_golang/internal/services"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:password@localhost:5432/bankdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewAccountRepository(db)
	service := services.NewAccountService(repo)
	handler := handlers.NewAccountHandler(service)

	http.HandleFunc("/accounts", handler.CreateAccount)
	http.HandleFunc("/accounts/{id}", handler.GetAccountByID)
	http.HandleFunc("/accounts/deposit", handler.AddMoneyToBalance)
	http.HandleFunc("/accounts/transaction", handler.Transaction)

	log.Println("Server started on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
