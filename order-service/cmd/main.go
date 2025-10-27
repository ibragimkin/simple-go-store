package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"order-service/config"
	_ "order-service/docs"
	"order-service/internal/adapters/httphandler"
	"order-service/internal/application/service"
	"order-service/internal/infrastructure/kafka"
	"order-service/internal/infrastructure/postgres"
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
	orderDb, err := postgres.NewPgOrderDb(db)
	if err != nil {
		log.Fatalf("failed to connect to account database: %v", err)
	}
	orderService := service.NewOrderService(orderDb)
	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaResponseTopic, cfg.KafkaGroupID)
	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaRequestTopic)
	messageBus := kafka.NewMessageBus(consumer, producer)
	go messageBus.StartReading(ctx)
	httpHandler := httphandler.NewOrderHandler(ctx, orderService, messageBus)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /orders/{id}", httpHandler.GetOrder)
	mux.HandleFunc("POST /orders", httpHandler.CreateOrder)
	mux.HandleFunc("GET /users/{id}/orders", httpHandler.GetUserOrders)
	mux.HandleFunc("PATCH /orders/{id}", httpHandler.PayOrder)
	mux.Handle("/swagger/order/", httpSwagger.WrapHandler)

	server := &http.Server{Addr: ":" + cfg.HttpPort, Handler: mux}
	log.Printf("Listening on port: %s", cfg.HttpPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
