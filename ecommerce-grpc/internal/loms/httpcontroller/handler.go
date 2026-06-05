package httpcontroller

import (
	"errors"
	"net/http"

	"grpc_pr/internal/httpx"
	"grpc_pr/internal/loms/entity"
	"grpc_pr/internal/loms/repository"
	"grpc_pr/internal/loms/service"
)

type Handler struct {
	orderService *service.OrderService
	stocksRepo   repository.StocksRepository
	products     map[uint32]entity.Product
}

func NewHandler(orderService *service.OrderService, stocksRepo repository.StocksRepository, products map[uint32]entity.Product) *Handler {
	return &Handler{orderService: orderService, stocksRepo: stocksRepo, products: products}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/products", h.listProducts)
	mux.HandleFunc("/product/list", h.listProducts)
	mux.HandleFunc("/product/info", h.getProductInfo)
	mux.HandleFunc("/product/create", h.createProduct)
	mux.HandleFunc("/product", h.getProductInfo)
	mux.HandleFunc("/stock/info", h.getStock)
	mux.HandleFunc("/stock/create", h.setStock)
	mux.HandleFunc("/stock", h.getStock)
	mux.HandleFunc("/order/create", h.createOrder)
	mux.HandleFunc("/order/info", h.getOrder)
	mux.HandleFunc("/order", h.getOrder)
	mux.HandleFunc("/order/pay", h.payOrder)
	mux.HandleFunc("/order/cancel", h.cancelOrder)

	mux.HandleFunc("/v1/products", h.listProducts)
	mux.HandleFunc("/v1/product/list", h.listProducts)
	mux.HandleFunc("/v1/product/info", h.getProductInfo)
	mux.HandleFunc("/v1/product/create", h.createProduct)
	mux.HandleFunc("/v1/stock/info", h.getStock)
	mux.HandleFunc("/v1/stock/create", h.setStock)
	mux.HandleFunc("/v1/order/create", h.createOrder)
	mux.HandleFunc("/v1/order/info", h.getOrder)
	mux.HandleFunc("/v1/order/pay", h.payOrder)
	mux.HandleFunc("/v1/order/cancel", h.cancelOrder)
	mux.HandleFunc("/v1/stock/set", h.setStock)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "loms"})
}

func mapLOMSHTTPError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrProductNotFound), errors.Is(err, entity.ErrOrderNotFound):
		httpx.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, entity.ErrInsufficientStock), errors.Is(err, entity.ErrOrderAlreadyPaid), errors.Is(err, entity.ErrOrderAlreadyCancelled), errors.Is(err, entity.ErrInvalidOrderStatus):
		httpx.WriteError(w, http.StatusPreconditionFailed, err.Error())
	case errors.Is(err, entity.ErrEmptyOrderItems):
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
	default:
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}
