package payment

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/payment"
	"shopnexus-go-service/internal/transport/connect/handler/common"
	"shopnexus-go-service/internal/transport/connect/interceptor/auth"
	"shopnexus-go-service/internal/utils/ptr"

	"connectrpc.com/connect"
	paymentv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1"
)

// GetRefund implements the GetRefund method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) GetRefund(ctx context.Context, req *connect.Request[paymentv1.GetRefundRequest]) (*connect.Response[paymentv1.GetRefundResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	refund, err := s.service.GetRefund(ctx, payment.GetRefundParams{
		UserID:   claims.UserID,
		RefundID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.GetRefundResponse{
		Data: convertRefundToProto(refund),
	}), nil
}

// ListRefunds implements the ListRefunds method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) ListRefunds(ctx context.Context, req *connect.Request[paymentv1.ListRefundsRequest]) (*connect.Response[paymentv1.ListRefundsResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	var (
		method *model.RefundMethod
		status *model.Status
	)
	if req.Msg.Method != nil {
		method = ptr.ToPtr(convertRefundMethod(*req.Msg.Method))
	}
	if req.Msg.Status != nil {
		status = ptr.ToPtr(convertStatus(*req.Msg.Status))
	}

	refunds, err := s.service.ListRefunds(ctx, payment.ListRefundsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		AccountID:          claims.UserID,
		Role:               claims.Role,
		ProductOnPaymentID: req.Msg.ProductOnPaymentId,
		Method:             method,
		Status:             status,
		Reason:             req.Msg.Reason,
		Address:            req.Msg.Address,
		DateCreatedFrom:    req.Msg.DateCreatedFrom,
		DateCreatedTo:      req.Msg.DateCreatedTo,
	})
	if err != nil {
		return nil, err
	}

	protoRefunds := make([]*paymentv1.Refund, 0, len(refunds.Data))
	for _, r := range refunds.Data {
		protoRefunds = append(protoRefunds, convertRefundToProto(r))
	}

	return connect.NewResponse(&paymentv1.ListRefundsResponse{
		Data:       protoRefunds,
		Pagination: common.ToProtoPaginationResponse(refunds),
	}), nil
}

// CreateRefund implements the CreateRefund method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) CreateRefund(ctx context.Context, req *connect.Request[paymentv1.CreateRefundRequest]) (*connect.Response[paymentv1.CreateRefundResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	refund, err := s.service.CreateRefund(ctx, payment.CreateRefundParams{
		UserID:             claims.UserID,
		ProductOnPaymentID: req.Msg.ProductOnPaymentId,
		Method:             convertRefundMethod(req.Msg.Method),
		Reason:             req.Msg.Reason,
		Address:            req.Msg.Address,
		Resources:          req.Msg.Resources,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CreateRefundResponse{
		Id: refund.ID,
	}), nil
}

// UpdateRefund implements the UpdateRefund method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) UpdateRefund(ctx context.Context, req *connect.Request[paymentv1.UpdateRefundRequest]) (*connect.Response[paymentv1.UpdateRefundResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	var (
		method    *model.RefundMethod
		status    *model.Status
		resources *[]string
	)
	if req.Msg.Method != nil {
		method = ptr.ToPtr(convertRefundMethod(*req.Msg.Method))
	}
	if req.Msg.Status != nil {
		status = ptr.ToPtr(convertStatus(*req.Msg.Status))
	}
	if req.Msg.Resources != nil {
		resources = &req.Msg.Resources
	}

	err = s.service.UpdateRefund(ctx, payment.UpdateRefundParams{
		ID:        req.Msg.Id,
		Role:      claims.Role,
		UserID:    claims.UserID,
		Method:    method,
		Status:    status,
		Reason:    req.Msg.Reason,
		Address:   req.Msg.Address,
		Resources: resources,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.UpdateRefundResponse{
		Id: req.Msg.Id,
	}), nil
}

// CancelRefund implements the CancelRefund method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) CancelRefund(ctx context.Context, req *connect.Request[paymentv1.CancelRefundRequest]) (*connect.Response[paymentv1.CancelRefundResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.CancelRefund(ctx, payment.CancelRefundParams{
		UserID:   claims.UserID,
		RefundID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CancelRefundResponse{}), nil
}
