package controller

import (
	"errors"

	"grpc_pr/internal/loms/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapLOMSError(err error) error {
	switch {
	case errors.Is(err, entity.ErrProductNotFound), errors.Is(err, entity.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entity.ErrInsufficientStock),
		errors.Is(err, entity.ErrOrderAlreadyPaid),
		errors.Is(err, entity.ErrOrderAlreadyCancelled),
		errors.Is(err, entity.ErrInvalidOrderStatus):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, entity.ErrEmptyOrderItems):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
