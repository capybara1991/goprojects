package httpcontroller

import (
	"encoding/json"
	"net/http"

	"grpc_pr/internal/httpx"
)

func (h *Handler) getStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sku, ok, err := httpx.Uint32Query(r, "sku")
	if err != nil || !ok || sku == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "sku must be positive")
		return
	}

	count, err := h.stocksRepo.GetStocks(r.Context(), sku)
	if err != nil {
		mapLOMSHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, stockResponse{SKU: sku, Count: count})
}

func (h *Handler) setStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req stockCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	sku := req.sku()
	if sku == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "sku must be positive")
		return
	}

	count := req.count()
	if err := h.stocksRepo.SetStocks(r.Context(), sku, count); err != nil {
		mapLOMSHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, stockResponse{SKU: sku, Count: count})
}
