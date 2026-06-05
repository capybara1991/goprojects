package httpcontroller

import "grpc_pr/internal/loms/entity"

type orderItemRequest struct {
	SKU   uint32 `json:"sku"`
	Count uint32 `json:"count"`
}

type createOrderRequest struct {
	UserID      int64              `json:"user_id"`
	UserIDCamel int64              `json:"userId"`
	Items       []orderItemRequest `json:"items"`
	Item        []orderItemRequest `json:"item"`
}

func (r createOrderRequest) userID() int64 {
	if r.UserID != 0 {
		return r.UserID
	}
	return r.UserIDCamel
}

func (r createOrderRequest) orderItems() []orderItemRequest {
	if len(r.Items) > 0 {
		return r.Items
	}
	return r.Item
}

type orderIDRequest struct {
	OrderID      uint32 `json:"order_id"`
	OrderIDCamel uint32 `json:"orderId"`
}

func (r orderIDRequest) orderID() uint32 {
	if r.OrderID != 0 {
		return r.OrderID
	}
	return r.OrderIDCamel
}

type productCreateRequest struct {
	SKU        uint32 `json:"sku"`
	SKUCamel   uint32 `json:"SKU"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Price      uint64 `json:"price"`
	Stock      uint64 `json:"stock"`
	Count      uint64 `json:"count"`
	StockCount uint64 `json:"stock_count"`
	StockCamel uint64 `json:"stockCount"`
}

func (r productCreateRequest) sku() uint32 {
	if r.SKU != 0 {
		return r.SKU
	}
	return r.SKUCamel
}

func (r productCreateRequest) name() string {
	if r.Name != "" {
		return r.Name
	}
	return r.Title
}

func (r productCreateRequest) price() uint64 {
	return r.Price
}

func (r productCreateRequest) stock() uint64 {
	if r.Stock != 0 {
		return r.Stock
	}
	if r.Count != 0 {
		return r.Count
	}
	if r.StockCount != 0 {
		return r.StockCount
	}
	return r.StockCamel
}

type stockCreateRequest struct {
	SKU        uint32 `json:"sku"`
	SKUCamel   uint32 `json:"SKU"`
	Count      uint64 `json:"count"`
	Stock      uint64 `json:"stock"`
	StockCount uint64 `json:"stock_count"`
	StockCamel uint64 `json:"stockCount"`
}

func (r stockCreateRequest) sku() uint32 {
	if r.SKU != 0 {
		return r.SKU
	}
	return r.SKUCamel
}

func (r stockCreateRequest) count() uint64 {
	if r.Count != 0 {
		return r.Count
	}
	if r.Stock != 0 {
		return r.Stock
	}
	if r.StockCount != 0 {
		return r.StockCount
	}
	return r.StockCamel
}

type productResponse struct {
	SKU   uint32 `json:"sku"`
	Name  string `json:"name"`
	Price uint64 `json:"price"`
}

type stockResponse struct {
	SKU   uint32 `json:"sku"`
	Count uint64 `json:"count"`
}

type orderItemResponse struct {
	SKU   uint32 `json:"sku"`
	Count uint32 `json:"count"`
}

type orderResponse struct {
	OrderID      uint32              `json:"order_id"`
	OrderIDCamel uint32              `json:"orderId,omitempty"`
	UserID       int64               `json:"user_id"`
	UserIDCamel  int64               `json:"userId,omitempty"`
	Status       string              `json:"status"`
	Items        []orderItemResponse `json:"items"`
}

func statusToString(status entity.Status) string {
	switch status {
	case entity.StatusAwaitingPayment:
		return "AWAITING_PAYMENT"
	case entity.StatusPaid:
		return "PAID"
	case entity.StatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

func makeOrderResponse(order *entity.Order) orderResponse {
	items := make([]orderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, orderItemResponse{SKU: item.SKU, Count: item.Count})
	}
	return orderResponse{
		OrderID:      order.ID,
		OrderIDCamel: order.ID,
		UserID:       order.UserID,
		UserIDCamel:  order.UserID,
		Status:       statusToString(order.Status),
		Items:        items,
	}
}
