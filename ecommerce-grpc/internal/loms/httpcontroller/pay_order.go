package httpcontroller

import (
	"net/http"

	"grpc_pr/internal/httpx"
)

func (h *Handler) payOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req orderIDRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderID := req.orderID()
	if orderID == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "order_id must be positive")
		return
	}

	if err := h.orderService.PayOrder(r.Context(), orderID); err != nil {
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
