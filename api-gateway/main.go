package main

import (
	"api-gateway/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

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

	r.Route("/orders", func(r chi.Router) {
		r.Handle("/*", orderProxy)
	})
	r.Route("/users/{id}/orders", func(r chi.Router) {
		r.Handle("/*", orderProxy)
	})
	r.Route("/users/{id}/account", func(r chi.Router) {
		r.Handle("/*", orderProxy)
	})

	r.Route("/accounts", func(r chi.Router) {
		r.Handle("/*", paymentProxy)
	})

	r.Route("/swagger/payment", func(r chi.Router) {
		r.Handle("/*", paymentProxy)
	})

	r.Route("/swagger/order", func(r chi.Router) {
		r.Handle("/*", orderProxy)
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	log.Println("API Gateway started on :" + cfg.HttpPort)
	if err := http.ListenAndServe(":"+cfg.HttpPort, r); err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}

}
