// cmd/main.go

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/everest/bheri/config"
	"github.com/everest/bheri/database"
	"github.com/everest/bheri/modules/user/controller"
	userrepo "github.com/everest/bheri/modules/user/repository"
	userservice "github.com/everest/bheri/modules/user/service"
	walletcontroller "github.com/everest/bheri/modules/wallet/controller"
	walletrepo "github.com/everest/bheri/modules/wallet/repository"
	walletservice "github.com/everest/bheri/modules/wallet/service"
	"github.com/everest/bheri/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger := config.NewLogger()

	db, err := database.NewPostgres()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	userRepository := userrepo.NewUserRepository(db)
	userService := userservice.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	walletRepository := walletrepo.NewWalletRepository(db)
	walletService := walletservice.NewWalletService(walletRepository)
	walletController := walletcontroller.NewWalletController(walletService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/api/v1", func(r chi.Router) {

		r.Post("/auth/register", userController.Register)
		r.Post("/auth/login", userController.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)

			r.Get("/me", userController.Me)

			r.Post("/wallets", walletController.Create)
			r.Get("/wallets", walletController.List)
			r.Get("/wallets/{id}", walletController.Get)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Printf("BHERI API running on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Fatal(err)
	}
}
