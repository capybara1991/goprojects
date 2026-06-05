package httpcontroller

import (
	"errors"
	"net/http"

	"grpc_pr/internal/cart/entity"
	"grpc_pr/internal/httpx"
)

type listCartItemResponse struct {
	SKU   uint32 `json:"sku"`
	Count uint32 `json:"count"`
	Name  string `json:"name"`
	Price uint64 `json:"price"`
}

type listCartResponse struct {
	Items           []listCartItemResponse `json:"items"`
	TotalPrice      uint64                 `json:"total_price"`
	TotalPriceCamel uint64                 `json:"totalPrice,omitempty"`
}

func (h *Handler) listCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var userID int64
	if r.Method == http.MethodGet {
		parsed, ok, err := httpx.Int64Query(r, "user_id", "userId")
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		if ok {
			userID = parsed
		}
	} else {
		var req userRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		userID = req.userID()
	}

	if userID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "user_id must be positive")
		return
	}

	items, total, err := h.itemService.ListCart(r.Context(), userID)
	if err != nil {
		if errors.Is(err, entity.ErrCartIsEmpty) {
			httpx.WriteJSON(w, http.StatusOK, listCartResponse{Items: []listCartItemResponse{}, TotalPrice: 0})
			return
		}
		mapCartHTTPError(w, err)
		return
	}

	respItems := make([]listCartItemResponse, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, listCartItemResponse{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	httpx.WriteJSON(w, http.StatusOK, listCartResponse{Items: respItems, TotalPrice: total, TotalPriceCamel: total})
}
