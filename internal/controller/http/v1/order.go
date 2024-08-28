package v1

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/internal/repo/cache"
	"github.com/v7ktory/wb_task_one/internal/repo/pgdb"
)

type orderRouter struct {
	cache     cache.Cache[string, *entity.Order]
	orderRepo pgdb.Order
	logger    *slog.Logger
}

func newOrderRouter(cache cache.Cache[string, *entity.Order], orderRepo pgdb.Order, logger *slog.Logger) http.Handler {
	o := &orderRouter{
		cache:     cache,
		orderRepo: orderRepo,
		logger:    logger,
	}
	mux := http.NewServeMux()
	mux.Handle("GET /order/", o.orderHomeHandler())
	mux.Handle("GET /order/my/{uid}", o.getOrderHandler())
	mux.Handle("GET /order/health", o.checkHealthHandler())

	return mux
}

func (o *orderRouter) orderHomeHandler() http.HandlerFunc {
	const op = "http.order.go - orderHomeHandler"

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./ui/templates/main.html")
		if err != nil {
			o.logger.Error("Error parsing template", slog.Any("error", err.Error()), slog.Any("operation", op))
			encode(w, http.StatusInternalServerError, "Error parsing template")
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			o.logger.Error("Error executing template", slog.Any("error", err.Error()), slog.Any("operation", op))
			encode(w, http.StatusInternalServerError, "Error executing template")
		}
	}
}
func (o *orderRouter) getOrderHandler() http.HandlerFunc {
	const op = "http.order.go - getOrderHandler"

	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.PathValue("uid")
		order, ok := o.cache.Get(uid)
		if !ok || order == nil {
			o.logger.Error("Order not found", slog.Any("uid", uid), slog.Any("operation", op))
			tmpl, err := template.ParseFiles("./ui/templates/not_found.html")
			if err != nil {
				o.logger.Error("Error parsing template", slog.Any("error", err.Error()), slog.Any("operation", op))
				encode(w, http.StatusInternalServerError, "Error parsing template")
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		err := o.orderRepo.UpdateOrderTime(r.Context(), uid)
		if err != nil {
			o.logger.Error("Error updating order time", slog.Any("error", err.Error()), slog.Any("operation", op))
			encode(w, http.StatusInternalServerError, "Error updating order time")
			return
		}

		tmpl, err := template.ParseFiles("./ui/templates/order.html")
		if err != nil {
			o.logger.Error("Error parsing template", slog.Any("error", err.Error()), slog.Any("operation", op))
			encode(w, http.StatusInternalServerError, "Error parsing template")
			return
		}

		tmpl.Execute(w, order)
	}
}
func (o *orderRouter) checkHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o.logger.Debug("healz")
		encode(w, http.StatusOK, "OK")
	}
}
