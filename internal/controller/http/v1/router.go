package v1

import (
	"log/slog"
	"net/http"

	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/internal/repo/cache"
)

func AddRoutes(mux *http.ServeMux, logger *slog.Logger, cache cache.CacheRepo[string, *entity.Order]) {
	// Handle Css files
	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle API routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", newOrderRouter(logger, cache)))
}
