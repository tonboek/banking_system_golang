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

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/login", userHandler.Login)


	accRepo := repository.NewAccountRepository(db)
	accService := services.NewAccountService(accRepo)
	accHandler := handlers.NewAccountHandler(accService)

	http.HandleFunc("/accounts", accHandler.CreateAccount)
	http.HandleFunc("/accounts/{id}", accHandler.GetAccountByID)
	http.HandleFunc("/accounts/deposit", accHandler.AddMoneyToBalance)
	http.HandleFunc("/accounts/transaction", accHandler.Transaction)

	log.Println("Server started on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
