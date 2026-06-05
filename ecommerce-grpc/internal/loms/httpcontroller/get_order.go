package httpcontroller

import (
	"net/http"

	"grpc_pr/internal/httpx"
)

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	orderID, ok, err := httpx.Uint32Query(r, "order_id", "orderId")
	if err != nil || !ok || orderID == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "order_id must be positive")
		return
	}

	order, err := h.orderService.GetOrder(r.Context(), orderID)
	if err != nil {
		mapLOMSHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, makeOrderResponse(order))
}
