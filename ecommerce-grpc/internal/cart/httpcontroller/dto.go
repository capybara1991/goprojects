package httpcontroller

type userRequest struct {
	UserID      int64 `json:"user_id"`
	UserIDCamel int64 `json:"userId"`
}

func (r userRequest) userID() int64 {
	if r.UserID != 0 {
		return r.UserID
	}
	return r.UserIDCamel
}

type addItemRequest struct {
	UserID      int64  `json:"user_id"`
	UserIDCamel int64  `json:"userId"`
	SKU         uint32 `json:"sku"`
	Count       uint32 `json:"count"`
}

func (r addItemRequest) userID() int64 {
	if r.UserID != 0 {
		return r.UserID
	}
	return r.UserIDCamel
}

type deleteItemRequest struct {
	UserID      int64  `json:"user_id"`
	UserIDCamel int64  `json:"userId"`
	SKU         uint32 `json:"sku"`
}

func (r deleteItemRequest) userID() int64 {
	if r.UserID != 0 {
		return r.UserID
	}
	return r.UserIDCamel
}
