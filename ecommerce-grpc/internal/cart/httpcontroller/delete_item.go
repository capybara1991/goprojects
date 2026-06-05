package httpcontroller

import (
	"net/http"

	"grpc_pr/internal/httpx"
)

func (h *Handler) deleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req deleteItemRequest
	if r.Method == http.MethodDelete && r.URL.Query().Get("user_id") != "" {
		userID, _, err := httpx.Int64Query(r, "user_id", "userId")
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		sku, _, err := httpx.Uint32Query(r, "sku")
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid sku")
			return
		}
		req.UserID = userID
		req.SKU = sku
	} else if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := req.userID()
	if userID <= 0 || req.SKU == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "user_id and sku must be positive")
		return
	}

	if err := h.itemService.DeleteItem(r.Context(), userID, req.SKU); err != nil {
		mapCartHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "DELETED"})
}
