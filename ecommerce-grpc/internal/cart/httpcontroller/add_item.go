package httpcontroller

import (
	"net/http"

	"grpc_pr/internal/httpx"
)

func (h *Handler) addItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req addItemRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := req.userID()
	if userID <= 0 || req.SKU == 0 || req.Count == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "user_id, sku and count must be positive")
		return
	}

	if err := h.itemService.AddItem(r.Context(), userID, req.SKU, req.Count); err != nil {
		mapCartHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "item added to cart"})
}
