// cmd/main.go

package main

import (
	"context"
	"net/http"
	"os"

	"github.com/everest/bheri/config"
	"github.com/everest/bheri/database"

	usercontroller "github.com/everest/bheri/modules/user/controller"
	userrepo "github.com/everest/bheri/modules/user/repository"
	userservice "github.com/everest/bheri/modules/user/service"

	walletcontroller "github.com/everest/bheri/modules/wallet/controller"
	walletrepo "github.com/everest/bheri/modules/wallet/repository"
	walletservice "github.com/everest/bheri/modules/wallet/service"

	agentcontroller "github.com/everest/bheri/modules/agent/controller"
	agentrepo "github.com/everest/bheri/modules/agent/repository"
	agentservice "github.com/everest/bheri/modules/agent/service"

	policycontroller "github.com/everest/bheri/modules/policy/controller"
	policyrepo "github.com/everest/bheri/modules/policy/repository"
	policyservice "github.com/everest/bheri/modules/policy/service"

	sessioncontroller "github.com/everest/bheri/modules/session/controller"
	sessionrepo "github.com/everest/bheri/modules/session/repository"
	sessionservice "github.com/everest/bheri/modules/session/service"

	paymentcontroller "github.com/everest/bheri/modules/payment/controller"
	paymentrepo "github.com/everest/bheri/modules/payment/repository"
	paymentservice "github.com/everest/bheri/modules/payment/service"

	apiservicecontroller "github.com/everest/bheri/modules/api_service/controller"
	apiservicerepo "github.com/everest/bheri/modules/api_service/repository"
	apiserviceservice "github.com/everest/bheri/modules/api_service/service"

	dashboardcontroller "github.com/everest/bheri/modules/dashboard/controller"
	dashboardrepo "github.com/everest/bheri/modules/dashboard/repository"
	dashboardservice "github.com/everest/bheri/modules/dashboard/service"

	"github.com/everest/bheri/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func main() {

	logger := config.NewLogger()

	db, err := database.NewPostgres()
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close(context.Background())

	// =========================================================
	// USER
	// =========================================================

	userRepository :=
		userrepo.NewUserRepository(db)

	userService :=
		userservice.NewUserService(
			userRepository,
		)

	userController :=
		usercontroller.NewUserController(
			userService,
		)

		//dahboard
		dashboardRepository :=
	dashboardrepo.NewDashboardRepository(db)

dashboardService :=
	dashboardservice.NewDashboardService(
		dashboardRepository,
	)

dashboardController :=
	dashboardcontroller.NewDashboardController(
		dashboardService,
	)

	// =========================================================
	// WALLET
	// =========================================================

	walletRepository :=
		walletrepo.NewWalletRepository(db)

	walletService :=
		walletservice.NewWalletService(
			walletRepository,
		)

	walletController :=
		walletcontroller.NewWalletController(
			walletService,
		)

	// =========================================================
	// AGENT
	// =========================================================

	agentRepository :=
		agentrepo.NewAgentRepository(db)

	agentService :=
		agentservice.NewAgentService(
			agentRepository,
		)

	agentController :=
		agentcontroller.NewAgentController(
			agentService,
		)

	// =========================================================
	// POLICY
	// =========================================================

	policyRepository :=
		policyrepo.NewPolicyRepository(db)

	policyService :=
		policyservice.NewPolicyService(
			policyRepository,
		)

	policyController :=
		policycontroller.NewPolicyController(
			policyService,
		)

	// =========================================================
	// SESSION
	// =========================================================

	sessionRepository :=
		sessionrepo.NewSessionRepository(db)

	sessionService :=
		sessionservice.NewSessionService(
			sessionRepository,
		)

	sessionController :=
		sessioncontroller.NewSessionController(
			sessionService,
		)

	// =========================================================
	// API SERVICES
	// =========================================================

	apiServiceRepository :=
		apiservicerepo.NewAPIServiceRepository(db)

	apiServiceService :=
		apiserviceservice.NewAPIServiceService(
			apiServiceRepository,
		)

	apiServiceController :=
		apiservicecontroller.NewAPIServiceController(
			apiServiceService,
		)

	// =========================================================
	// PAYMENT
	// =========================================================

	paymentRepository :=
		paymentrepo.NewPaymentRepository(db)

	/*
		Payment service requires:

		- PaymentRepository
		- PolicyEngine
		- SessionService

		Your policyService must implement PolicyEngine.
	*/

	paymentService :=
		paymentservice.NewPaymentService(
			paymentRepository,
			policyService,
			sessionService,
		)

	paymentController :=
		paymentcontroller.NewPaymentController(
			paymentService,
		)

	// =========================================================
	// ROUTER
	// =========================================================

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/api/v1", func(r chi.Router) {

		// =====================================================
		// AUTH
		// =====================================================

		r.Post(
			"/auth/register",
			userController.Register,
		)

		r.Post(
			"/auth/login",
			userController.Login,
		)

		// =====================================================
		// USER
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			r.Get(
				"/me",
				userController.Me,
			)
		})

		// =====================================================
		// WALLET
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			r.Post(
				"/wallets",
				walletController.Create,
			)

			r.Get(
				"/wallets",
				walletController.List,
			)

			r.Get(
				"/wallets/{id}",
				walletController.Get,
			)
		})

		// =====================================================
		// AGENTS
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			r.Post(
				"/agents",
				agentController.Create,
			)

			r.Get(
				"/agents",
				agentController.List,
			)

			r.Get(
				"/agents/{id}",
				agentController.Get,
			)

			r.Put(
				"/agents/{id}",
				agentController.Update,
			)

			r.Delete(
				"/agents/{id}",
				agentController.Delete,
			)
		})

		// =====================================================
		// POLICIES
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			r.Post(
				"/policies",
				policyController.Create,
			)

			r.Get(
				"/agents/{agentID}/policy",
				policyController.Get,
			)

			r.Post(
				"/agents/{agentID}/policy/evaluate",
				policyController.Evaluate,
			)
		})
		//=====================================================
		// Dashboard
		// =====================================================
r.Group(func(r chi.Router) {

	r.Use(middleware.Auth)

	r.Get(
		"/dashboard",
		dashboardController.Get,
	)
})

		// =====================================================
		// SESSIONS
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			r.Post(
				"/agents/{agentID}/sessions",
				sessionController.Create,
			)

			r.Get(
				"/agents/{agentID}/sessions",
				sessionController.List,
			)

			r.Post(
				"/sessions/{id}/revoke",
				sessionController.Revoke,
			)
		})

		// =====================================================
		// PAYMENTS
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			r.Post(
				"/payments",
				paymentController.Create,
			)

			r.Get(
				"/payments/{id}",
				paymentController.Get,
			)
		})

		// =====================================================
		// API SERVICES
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			r.Post(
				"/api-services",
				apiServiceController.Create,
			)

			r.Get(
				"/api-services",
				apiServiceController.List,
			)

			r.Get(
				"/api-services/active",
				apiServiceController.ListActive,
			)

			r.Get(
				"/api-services/{id}",
				apiServiceController.Get,
			)

			r.Put(
				"/api-services/{id}",
				apiServiceController.Update,
			)

			r.Delete(
				"/api-services/{id}",
				apiServiceController.Delete,
			)
		})
	})

	// =========================================================
	// SERVER
	// =========================================================

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	logger.Printf(
		"BHERI API running on :%s",
		port,
	)

	if err := http.ListenAndServe(
		":"+port,
		router,
	); err != nil {
		logger.Fatal(err)
	}
}
