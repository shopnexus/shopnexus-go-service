package payment

import (
	"context"
	"shopnexus-go-service/internal/grpc/server/interceptor"
	"shopnexus-go-service/internal/model"
	"strconv"

	"connectrpc.com/connect"
	paymentv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1"
)

// GetRefund implements the GetRefund method from PaymentServiceHandler
func (s *PaymentServer) GetRefund(ctx context.Context, req *connect.Request[paymentv1.GetRefundRequest]) (*connect.Response[paymentv1.GetRefundResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	refundID, err := strconv.ParseInt(req.Msg.RefundId, 10, 64)
	if err != nil {
		return nil, err
	}

	refund, err := s.service.GetRefund(ctx, refundID, claims.UserID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.GetRefundResponse{
		Refund: convertRefundToProto(refund),
	}), nil
}

// ListRefunds implements the ListRefunds method from PaymentServiceHandler
func (s *PaymentServer) ListRefunds(ctx context.Context, req *connect.Request[paymentv1.ListRefundsRequest]) (*connect.Response[paymentv1.ListRefundsResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	refunds, err := s.service.ListRefunds(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	protoRefunds := make([]*paymentv1.Refund, 0, len(refunds))
	for _, r := range refunds {
		protoRefunds = append(protoRefunds, convertRefundToProto(r))
	}

	return connect.NewResponse(&paymentv1.ListRefundsResponse{
		Refunds: protoRefunds,
	}), nil
}

// CreateRefund implements the CreateRefund method from PaymentServiceHandler
func (s *PaymentServer) CreateRefund(ctx context.Context, req *connect.Request[paymentv1.CreateRefundRequest]) (*connect.Response[paymentv1.CreateRefundResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	paymentID, err := strconv.ParseInt(req.Msg.PaymentId, 10, 64)
	if err != nil {
		return nil, err
	}

	refund, err := s.service.CreateRefund(ctx, payment.CreateRefundParams{
		UserID:    claims.UserID,
		PaymentID: paymentID,
		Method:    convertRefundMethod(req.Msg.Method),
		Reason:    req.Msg.Reason,
		Address:   req.Msg.Address,
		Resources: req.Msg.Resources,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CreateRefundResponse{
		RefundId: strconv.FormatInt(refund.ID, 10),
	}), nil
}

// UpdateRefund implements the UpdateRefund method from PaymentServiceHandler
func (s *PaymentServer) UpdateRefund(ctx context.Context, req *connect.Request[paymentv1.UpdateRefundRequest]) (*connect.Response[paymentv1.UpdateRefundResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	refundID, err := strconv.ParseInt(req.Msg.RefundId, 10, 64)
	if err != nil {
		return nil, err
	}

	refund, err := s.service.UpdateRefund(ctx, payment.UpdateRefundParams{
		ID:     refundID,
		UserID: claims.UserID,
		Status: convertStatus(req.Msg.Status),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.UpdateRefundResponse{
		Refund: convertRefundToProto(refund),
	}), nil
}

// CancelRefund implements the CancelRefund method from PaymentServiceHandler
func (s *PaymentServer) CancelRefund(ctx context.Context, req *connect.Request[paymentv1.CancelRefundRequest]) (*connect.Response[paymentv1.CancelRefundResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	refundID, err := strconv.ParseInt(req.Msg.RefundId, 10, 64)
	if err != nil {
		return nil, err
	}

	err = s.service.CancelRefund(ctx, refundID, claims.UserID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CancelRefundResponse{}), nil
}
