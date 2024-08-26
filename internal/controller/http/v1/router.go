package v1

import (
	"log/slog"
	"net/http"

	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/internal/repo/cache"
	"github.com/v7ktory/wb_task_one/internal/repo/pgdb"
)

func AddRoutes(mux *http.ServeMux, cache cache.CacheRepo[string, *entity.Order], pgRepo *pgdb.PgRepo, logger *slog.Logger) {
	// Handle Css files
	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle API routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", newOrderRouter(cache, pgRepo, logger)))
}
