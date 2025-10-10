package main

import (
	"api-gateway/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// создаём proxy для каждого микросервиса
func newReverseProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("invalid target URL: %v", err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	r := chi.NewRouter()

	// Прокси для микросервисов
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	orderProxy := newReverseProxy(cfg.OrderServiceURL)
	paymentProxy := newReverseProxy(cfg.PaymentServiceURL)

	// Группы маршрутов — всё, что начинается с /orders или /users/{id}/orders → order-service
	r.Route("/orders", func(r chi.Router) {
		r.Handle("/*", orderProxy) // проксирует всё, включая /orders/{id}
	})
	r.Route("/users/{id}/orders", func(r chi.Router) {
		r.Handle("/*", orderProxy)
	})

	// Всё, что начинается с /accounts → payment-service
	r.Route("/accounts", func(r chi.Router) {
		r.Handle("/*", paymentProxy)
	})

	r.Route("/swagger/payment", func(r chi.Router) {
		r.Handle("/*", paymentProxy)
	})

	// /swagger/order/* → order-service/swagger/*
	r.Route("/swagger/order", func(r chi.Router) {
		r.Handle("/*", orderProxy)
	})

	// Можно добавить health-check, чтобы удобно проверять доступность гейтвея
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	log.Println("API Gateway started on :" + cfg.HttpPort)
	if err := http.ListenAndServe(":"+cfg.HttpPort, r); err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}

}
