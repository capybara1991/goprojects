package httpcontroller

import (
	"encoding/json"
	"grpc_pr/internal/loms/entity"
	"net/http"
	"sort"
	"strconv"

	"grpc_pr/internal/httpx"
)

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	skus := make([]int, 0, len(h.products))
	for sku := range h.products {
		skus = append(skus, int(sku))
	}
	sort.Ints(skus)

	resp := make([]productResponse, 0, len(skus))
	for _, rawSKU := range skus {
		sku := uint32(rawSKU)
		product := h.products[sku]
		resp = append(resp, productResponse{SKU: sku, Name: product.Name, Price: product.Price})
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) getProductInfo(w http.ResponseWriter, r *http.Request) {
	var sku uint32

	switch r.Method {
	case http.MethodGet:
		rawSKU := r.URL.Query().Get("sku")
		if rawSKU == "" {
			httpx.WriteError(w, http.StatusBadRequest, "sku is required")
			return
		}

		parsedSKU, err := strconv.ParseUint(rawSKU, 10, 32)
		if err != nil || parsedSKU == 0 {
			httpx.WriteError(w, http.StatusBadRequest, "sku must be positive")
			return
		}

		sku = uint32(parsedSKU)

	case http.MethodPost:
		var req struct {
			SKU uint32 `json:"sku"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		if req.SKU == 0 {
			httpx.WriteError(w, http.StatusBadRequest, "sku must be positive")
			return
		}

		sku = req.SKU

	default:
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	product, ok := h.products[sku]
	if !ok {
		httpx.WriteError(w, http.StatusNotFound, "product not found")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, productResponse{
		SKU:   sku,
		Name:  product.Name,
		Price: product.Price,
	})
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req productCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	name := req.name()
	if name == "" {
		httpx.WriteError(w, http.StatusBadRequest, "product name is required")
		return
	}

	price := req.price()
	if price == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "price must be positive")
		return
	}

	sku := req.sku()

	if sku == 0 {
		var maxSKU uint32
		for existingSKU := range h.products {
			if existingSKU > maxSKU {
				maxSKU = existingSKU
			}
		}
		sku = maxSKU + 1
	}

	h.products[sku] = entity.Product{
		Name:  name,
		Price: price,
	}

	stock := req.stock()
	if stock == 0 {
		stock = 100
	}

	if err := h.stocksRepo.SetStocks(r.Context(), sku, stock); err != nil {
		mapLOMSHTTPError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, productResponse{
		SKU:   sku,
		Name:  name,
		Price: price,
	})
}
