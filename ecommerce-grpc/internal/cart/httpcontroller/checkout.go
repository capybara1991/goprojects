package httpcontroller

import (
	"net/http"

	"grpc_pr/internal/httpx"
)

type checkoutResponse struct {
	OrderID      int64 `json:"order_id"`
	OrderIDCamel int64 `json:"orderId,omitempty"`
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req userRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := req.userID()
	if userID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "user_id must be positive")
		return
	}

	orderID, err := h.itemService.Checkout(r.Context(), userID)
	if err != nil {
		mapCartHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, checkoutResponse{OrderID: orderID, OrderIDCamel: orderID})
}
