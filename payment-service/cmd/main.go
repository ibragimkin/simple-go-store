package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "payment-service/docs"
	"payment-service/internal/adapters/httphandler"
	"payment-service/internal/adapters/kafkahandler"
	"payment-service/internal/application/service"
	"payment-service/internal/config"
	"payment-service/internal/infrastructure/kafka"
	"payment-service/internal/infrastructure/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	var db *pgxpool.Pool
	db, err = pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	err = db.Ping(ctx)
	if err != nil {
		db.Close()
		log.Fatalf("Unable to ping database: %v", err)
	}
	defer db.Close()
	accountRepo, err := postgres.NewAccountDb(db)
	if err != nil {
		log.Fatalf("failed to connect to account database: %v", err)
	}
	transactionRepo, err := postgres.NewTransactionDb(db)
	if err != nil {
		log.Fatalf("failed to connect to transaciton database: %v", err)
	}
	accountService := service.NewAccountService(accountRepo)
	paymentService, err := service.NewPaymentService(accountRepo, transactionRepo)
	if err != nil {
		log.Fatalf("failed to initialize payment service: %v", err)
	}

	httpHandler := httphandler.NewAccountHandler(ctx, accountService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /accounts/{id}", httpHandler.GetAccount)
	mux.HandleFunc("PATCH /accounts/{id}", httpHandler.Deposit)
	mux.HandleFunc("POST /accounts", httpHandler.CreateAccount)
	mux.HandleFunc("GET /users/{id}/account", httpHandler.GetUsersAccount)
	mux.Handle("/swagger/payment/", httpSwagger.WrapHandler)
	messageBus := kafka.NewMessageBus(cfg.KafkaBrokers, cfg.KafkaConsumerTopic, cfg.KafkaProducerTopic, cfg.KafkaGroupID)
	kafkaHandler := kafkahandler.NewPaymentHandler(paymentService)
	go messageBus.Start(ctx, kafkaHandler)
	server := &http.Server{Addr: ":" + cfg.HttpPort, Handler: mux}
	log.Printf("Listening on port " + cfg.HttpPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
