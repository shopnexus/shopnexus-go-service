package server

import (
	"context"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/grpc/server/interceptor"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service"
	"strconv"
)

type RefundServer struct {
	pb.UnimplementedRefundServer
	service *service.RefundService
}

func NewRefundServer(service *service.RefundService) *RefundServer {
	return &RefundServer{service: service}
}

func (s *RefundServer) Create(ctx context.Context, req *pb.CreateRefundRequest) (*pb.CreateRefundResponse, error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	refund, err := s.service.Create(ctx, service.CreateRefundParams{
		UserID:    claims.UserID,
		PaymentID: req.PaymentId,
		Method:    convertRefundMethod(req.Method),
		Reason:    req.Reason,
		Address:   req.Address,
		// Resources: req.Resources,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateRefundResponse{
		RefundId: strconv.FormatInt(refund.ID, 10),
	}, nil
}

func convertRefundMethod(protoMethod pb.RefundMethod) model.RefundMethod {
	switch protoMethod {
	case pb.RefundMethod_DROP_OFF:
		return model.RefundMethodDropOff
	case pb.RefundMethod_PICK_UP:
		return model.RefundMethodPickUp
	case pb.RefundMethod_REFUND_METHOD_UNSPECIFIED:
		panic("refund method is unspecified")
	default:
		panic("unknown refund method")
	}
}
