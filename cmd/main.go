// cmd/main.go

package main

import (
	"net/http"
	"os"

	"github.com/everest/bheri/config"
	"github.com/everest/bheri/database"

	agentcontroller "github.com/everest/bheri/modules/agent/controller"
	agentrepo "github.com/everest/bheri/modules/agent/repository"
	agentservice "github.com/everest/bheri/modules/agent/service"

	apikeycontroller "github.com/everest/bheri/modules/api_key/controller"
	apikeyrepo "github.com/everest/bheri/modules/api_key/repository"
	apikeyservice "github.com/everest/bheri/modules/api_key/service"

	apiservicecontroller "github.com/everest/bheri/modules/api_service/controller"
	apiservicerepo "github.com/everest/bheri/modules/api_service/repository"
	apiserviceservice "github.com/everest/bheri/modules/api_service/service"

	blockchainclient "github.com/everest/bheri/modules/blockchain/client"
	blockchainservice "github.com/everest/bheri/modules/blockchain/service"

	dashboardcontroller "github.com/everest/bheri/modules/dashboard/controller"
	dashboardrepo "github.com/everest/bheri/modules/dashboard/repository"
	dashboardservice "github.com/everest/bheri/modules/dashboard/service"

	ledgercontroller "github.com/everest/bheri/modules/ledger/controller"
	ledgerrepo "github.com/everest/bheri/modules/ledger/repository"
	ledgerservice "github.com/everest/bheri/modules/ledger/service"

	paymentcontroller "github.com/everest/bheri/modules/payment/controller"
	paymentrepo "github.com/everest/bheri/modules/payment/repository"
	paymentservice "github.com/everest/bheri/modules/payment/service"

	policycontroller "github.com/everest/bheri/modules/policy/controller"
	policyrepo "github.com/everest/bheri/modules/policy/repository"
	policyservice "github.com/everest/bheri/modules/policy/service"

	sessioncontroller "github.com/everest/bheri/modules/session/controller"
	sessionrepo "github.com/everest/bheri/modules/session/repository"
	sessionservice "github.com/everest/bheri/modules/session/service"

	settlementcontroller "github.com/everest/bheri/modules/settlement/controller"
	settlementrepo "github.com/everest/bheri/modules/settlement/repository"
	settlementservice "github.com/everest/bheri/modules/settlement/service"

	transactioncontroller "github.com/everest/bheri/modules/transaction/controller"
	transactionrepo "github.com/everest/bheri/modules/transaction/repository"
	transactionservice "github.com/everest/bheri/modules/transaction/service"

	usercontroller "github.com/everest/bheri/modules/user/controller"
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

	// =========================================================
	// DATABASE
	// =========================================================

	db, err := database.NewPostgres()
	if err != nil {
		logger.Fatal(err)
	}


	// =========================================================
	// USER
	// =========================================================

	userRepository := userrepo.NewUserRepository(db)

	userService := userservice.NewUserService(
		userRepository,
	)

	userController := usercontroller.NewUserController(
		userService,
	)

	// =========================================================
	// WALLET
	// =========================================================

	walletRepository := walletrepo.NewWalletRepository(db)

	walletService := walletservice.NewWalletService(
		walletRepository,
	)

	walletController := walletcontroller.NewWalletController(
		walletService,
	)

	// =========================================================
	// AGENT
	// =========================================================

	agentRepository := agentrepo.NewAgentRepository(db)

	agentService := agentservice.NewAgentService(
		agentRepository,
	)

	agentController := agentcontroller.NewAgentController(
		agentService,
	)

	// =========================================================
	// POLICY
	// =========================================================

	policyRepository := policyrepo.NewPolicyRepository(db)

	policyService := policyservice.NewPolicyService(
		policyRepository,
	)

	policyController := policycontroller.NewPolicyController(
		policyService,
	)

	// =========================================================
	// SESSION
	// =========================================================

	sessionRepository := sessionrepo.NewSessionRepository(db)

	sessionService := sessionservice.NewSessionService(
		sessionRepository,
	)

	sessionController := sessioncontroller.NewSessionController(
		sessionService,
	)

	// =========================================================
	// API KEY
	// =========================================================

	apiKeyRepository := apikeyrepo.NewAPIKeyRepository(db)

	apiKeyService := apikeyservice.NewAPIKeyService(
		apiKeyRepository,
	)

	apiKeyController := apikeycontroller.NewAPIKeyController(
		apiKeyService,
	)

	// =========================================================
	// API SERVICES
	// =========================================================

	apiServiceRepository := apiservicerepo.NewAPIServiceRepository(db)

	apiServiceService := apiserviceservice.NewAPIServiceService(
		apiServiceRepository,
	)

	apiServiceController := apiservicecontroller.NewAPIServiceController(
		apiServiceService,
	)

	// =========================================================
	// BLOCKCHAIN
	// =========================================================

	solanaRPC := os.Getenv("SOLANA_RPC_URL")

	if solanaRPC == "" {
		solanaRPC = "https://api.devnet.solana.com"
	}

	solanaClient := blockchainclient.NewSolanaClient(
		solanaRPC,
	)

	blockchainService := blockchainservice.NewBlockchainService(
		solanaClient,
	)

	// =========================================================
	// PAYMENT
	// =========================================================

	paymentRepository := paymentrepo.NewPaymentRepository(db)

	paymentService := paymentservice.NewPaymentService(
		paymentRepository,
		policyService,
		sessionService,
	)

	paymentController := paymentcontroller.NewPaymentController(
		paymentService,
	)

	// =========================================================
	// SETTLEMENT
	// =========================================================

	settlementRepository := settlementrepo.NewSettlementRepository(db)

	settlementService := settlementservice.NewSettlementService(
		settlementRepository,
		blockchainService,
	)

	settlementController := settlementcontroller.NewSettlementController(
		settlementService,
	)

	// =========================================================
	// TRANSACTION
	// =========================================================

	transactionRepository := transactionrepo.NewTransactionRepository(db)

	transactionService := transactionservice.NewTransactionService(
		transactionRepository,
	)

	transactionController := transactioncontroller.NewTransactionController(
		transactionService,
	)

	// =========================================================
	// LEDGER
	// =========================================================

	ledgerRepository := ledgerrepo.NewLedgerRepository(db)

	ledgerService := ledgerservice.NewLedgerService(
		ledgerRepository,
	)

	ledgerController := ledgercontroller.NewLedgerController(
		ledgerService,
	)

	// =========================================================
	// DASHBOARD
	// =========================================================

	dashboardRepository := dashboardrepo.NewDashboardRepository(db)

	dashboardService := dashboardservice.NewDashboardService(
		dashboardRepository,
	)

	dashboardController := dashboardcontroller.NewDashboardController(
		dashboardService,
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
		// AUTHENTICATED ROUTES
		// =====================================================

		r.Group(func(r chi.Router) {

			r.Use(middleware.Auth)

			// =================================================
			// USER
			// =================================================

			r.Get(
				"/me",
				userController.Me,
			)

			// =================================================
			// DASHBOARD
			// =================================================

			r.Get(
				"/dashboard",
				dashboardController.Get,
			)

			// =================================================
			// WALLET
			// =================================================

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

			// =================================================
			// AGENTS
			// =================================================

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

			// =================================================
			// POLICIES
			// =================================================

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

			// =================================================
			// SESSIONS
			// =================================================

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

			// =================================================
			// API KEYS
			// =================================================

			r.Post(
				"/api-keys",
				apiKeyController.Create,
			)

			r.Get(
				"/api-keys",
				apiKeyController.List,
			)

			r.Get(
				"/api-keys/{id}",
				apiKeyController.Get,
			)

			r.Post(
				"/api-keys/{id}/revoke",
				apiKeyController.Revoke,
			)

			// =================================================
			// API SERVICES
			// =================================================

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

			// =================================================
			// PAYMENTS
			// =================================================

			r.Post(
				"/payments",
				paymentController.Create,
			)

			r.Get(
				"/payments/{id}",
				paymentController.Get,
			)

			// =================================================
			// TRANSACTIONS
			// =================================================

			r.Get(
				"/transactions",
				transactionController.List,
			)

			r.Get(
				"/transactions/{id}",
				transactionController.Get,
			)

			// =================================================
			// SETTLEMENTS
			// =================================================

			r.Post(
				"/settlements",
				settlementController.Create,
			)

			r.Get(
				"/settlements",
				settlementController.List,
			)

			r.Get(
				"/settlements/{id}",
				settlementController.Get,
			)

			r.Post(
				"/settlements/{id}/verify",
				settlementController.Verify,
			)

			// =================================================
			// LEDGER
			// =================================================

			r.Get(
				"/ledger",
				ledgerController.ListAccounts,
			)

			r.Get(
				"/ledger/{id}",
				ledgerController.GetAccount,
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
