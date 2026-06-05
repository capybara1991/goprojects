package httpcontroller

import (
	"errors"
	"net/http"

	"grpc_pr/internal/cart/entity"
	"grpc_pr/internal/cart/service"
	"grpc_pr/internal/httpx"
)

type Handler struct {
	itemService *service.ItemService
}

func NewHandler(itemService *service.ItemService) *Handler {
	return &Handler{itemService: itemService}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/cart/item/add", h.addItem)
	mux.HandleFunc("/cart/item", h.deleteItem)
	mux.HandleFunc("/cart/item/delete", h.deleteItem)
	mux.HandleFunc("/cart/list", h.listCart)
	mux.HandleFunc("/cart", h.listCart)
	mux.HandleFunc("/cart/checkout", h.checkout)
	mux.HandleFunc("/checkout", h.checkout)
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/v1/cart/item/add", h.addItem)
	mux.HandleFunc("/v1/cart/item/delete", h.deleteItem)
	mux.HandleFunc("/v1/cart/list", h.listCart)
	mux.HandleFunc("/v1/cart/checkout", h.checkout)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "cart"})
}

func mapCartHTTPError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrProductNotFound), errors.Is(err, entity.ErrCartNotFound), errors.Is(err, entity.ErrCartItemNotFound):
		httpx.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, entity.ErrInsufficientStock), errors.Is(err, entity.ErrCartIsEmpty):
		httpx.WriteError(w, http.StatusPreconditionFailed, err.Error())
	default:
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}
