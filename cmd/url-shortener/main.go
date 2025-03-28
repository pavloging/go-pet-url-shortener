package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// Get environment
	cfg := config.MustLoad()

	// Settings logger
	log := setupLogger(cfg.Env)
	log.Info("starting the project...", slog.String("env", cfg.Env), slog.String("version", "v1"))
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	// Settings and started database
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// Init router
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.RequestID) // Хороший middleware для логирования
	router.Use(mwLogger.New(log))    // Хороший middleware для логирования (custom)
	router.Use(middleware.Logger)    // Логирует все входящие запросы
	router.Use(middleware.Recoverer) // Перехватывает паники и возвращает 500
	router.Use(middleware.URLFormat) // Для красивых URL при подключении к обработчикам

	// Handlers with Auth
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.Auth.User: cfg.Auth.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	// Handlers without Auth
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	// Settings and started server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog() // Для локальный разработки - самописный logger
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
