package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Config struct {
	Log *zap.Logger
	DB  *pgxpool.Pool
}

func NewRouter(cfg Config) *chi.Mux {

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"Connect-Protocol-Vision",
		},
		ExposedHeaders:   []string{"Grpc-Status", "Grpc-Message", "Set-Cookie"},
		AllowCredentials: true,
		Debug:            true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

	cfg.Log.Info("Router initialized successfully")

	return r
}
