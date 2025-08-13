package generic

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GenericHandler cung cấp các method chung cho tất cả handler
type GenericHandler struct{}

// HandleOperation xử lý tất cả các operation với pattern chung
func HandleOperation[T any, R any](
	ctx context.Context,
	req T,
	operation func(context.Context, T) (R, error),
	operationName string,
) (R, error) {
	resp, err := operation(ctx, req)
	if err != nil {
		return *new(R), status.Errorf(codes.Internal, "failed to %s: %v", operationName, err)
	}
	return resp, nil
}

// HandleOperationWithID xử lý các operation cần ID parameter
func HandleOperationWithID[T any, R any, ID any](
	ctx context.Context,
	id ID,
	operation func(context.Context, ID) (R, error),
	operationName string,
) (R, error) {
	resp, err := operation(ctx, id)
	if err != nil {
		return *new(R), status.Errorf(codes.Internal, "failed to %s: %v", operationName, err)
	}
	return resp, nil
}

// HandleOperationWithString xử lý các operation cần string parameter
func HandleOperationWithString[T any, R any](
	ctx context.Context,
	str string,
	operation func(context.Context, string) (R, error),
	operationName string,
) (R, error) {
	resp, err := operation(ctx, str)
	if err != nil {
		return *new(R), status.Errorf(codes.Internal, "failed to %s: %v", operationName, err)
	}
	return resp, nil
}

// HandleListOperation xử lý các operation trả về list
func HandleListOperation[T any, R any](
	ctx context.Context,
	req T,
	operation func(context.Context, T) (R, error),
	operationName string,
) (R, error) {
	resp, err := operation(ctx, req)
	if err != nil {
		return *new(R), status.Errorf(codes.Internal, "failed to %s: %v", operationName, err)
	}
	return resp, nil
}

// NewGenericHandler tạo instance mới của GenericHandler
func NewGenericHandler() *GenericHandler {
	return &GenericHandler{}
}
