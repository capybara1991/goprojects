package controller

import (
	cartpb "grpc_pr/gen/cart"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateAddItemRequest(req *cartpb.AddItemRequest) error {
	if req.UserId <= 0 {
		return status.Error(codes.InvalidArgument, "user_id must be positive")
	}
	if req.Sku == 0 {
		return status.Error(codes.InvalidArgument, "sku must be positive")
	}
	if req.Count == 0 {
		return status.Error(codes.InvalidArgument, "count must be positive")
	}
	return nil
}

func validateDeleteItemRequest(req *cartpb.DeleteItemRequest) error {
	if req.UserId <= 0 {
		return status.Error(codes.InvalidArgument, "user_id must be positive")
	}
	if req.Sku == 0 {
		return status.Error(codes.InvalidArgument, "sku must be positive")
	}
	return nil
}

func validateUserID(userID int64) error {
	if userID <= 0 {
		return status.Error(codes.InvalidArgument, "user_id must be positive")
	}
	return nil
}
