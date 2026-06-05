package httpcontroller

import (
	"net/http"

	"grpc_pr/internal/httpx"
	"grpc_pr/internal/loms/entity"
)

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req createOrderRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := req.userID()
	if userID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "user_id must be positive")
		return
	}

	requestItems := req.orderItems()
	items := make([]entity.OrderItem, 0, len(requestItems))
	for _, item := range requestItems {
		if item.SKU == 0 || item.Count == 0 {
			httpx.WriteError(w, http.StatusBadRequest, "sku and count must be positive")
			return
		}
		items = append(items, entity.OrderItem{SKU: item.SKU, Count: item.Count})
	}

	orderID, err := h.orderService.CreateOrder(r.Context(), userID, items)
	if err != nil {
		mapLOMSHTTPError(w, err)
		return
	}

	order, err := h.orderService.GetOrder(r.Context(), orderID)
	if err != nil {
		mapLOMSHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, makeOrderResponse(order))
}
