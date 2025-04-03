package payment

import (
	"context"
	"net/http"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/server/connect/interceptor/auth"
	"shopnexus-go-service/internal/service/payment"
	"shopnexus-go-service/internal/util"

	common_grpc "shopnexus-go-service/internal/server/connect/handler/common"

	"connectrpc.com/connect"
	common_pb "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/common"
	paymentv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1/paymentv1connect"
)

type ImplementedPaymentServiceHandler struct {
	paymentv1connect.UnimplementedPaymentServiceHandler
	service *payment.PaymentService
}

func NewPaymentServiceHandler(paymentService *payment.PaymentService, opts ...connect.HandlerOption) (string, http.Handler) {
	return paymentv1connect.NewPaymentServiceHandler(&ImplementedPaymentServiceHandler{
		service: paymentService,
	}, opts...)
}

// GetPayment implements the GetPayment method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) GetPayment(ctx context.Context, req *connect.Request[paymentv1.GetPaymentRequest]) (*connect.Response[paymentv1.GetPaymentResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	payment, err := s.service.GetPayment(ctx, payment.GetPaymentParams{
		UserID:    claims.UserID,
		PaymentID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.GetPaymentResponse{
		Data: convertPaymentToProto(payment),
	}), nil
}

// ListPayments implements the ListPayments method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) ListPayments(ctx context.Context, req *connect.Request[paymentv1.ListPaymentsRequest]) (*connect.Response[paymentv1.ListPaymentsResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	var (
		method *model.PaymentMethod
		status *model.Status
	)

	if req.Msg.Method != nil {
		method = util.ToPtr(convertPaymentMethod(*req.Msg.Method))
	}

	if req.Msg.Status != nil {
		status = util.ToPtr(convertStatus(*req.Msg.Status))
	}

	payments, err := s.service.ListPayments(ctx, payment.ListPaymentsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		UserID:          &claims.UserID,
		Method:          method,
		Status:          status,
		Address:         req.Msg.Address,
		TotalFrom:       req.Msg.TotalFrom,
		TotalTo:         req.Msg.TotalTo,
		DateCreatedFrom: req.Msg.DateCreatedFrom,
		DateCreatedTo:   req.Msg.DateCreatedTo,
	})
	if err != nil {
		return nil, err
	}

	protoPayments := make([]*paymentv1.Payment, 0, len(payments.Data))
	for _, p := range payments.Data {
		protoPayments = append(protoPayments, convertPaymentToProto(p))
	}

	return connect.NewResponse(&paymentv1.ListPaymentsResponse{
		Data:       protoPayments,
		Pagination: common_grpc.ToProtoPaginationResponse(payments),
	}), nil
}

// CreatePayment implements the CreatePayment method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) CreatePayment(ctx context.Context, req *connect.Request[paymentv1.CreatePaymentRequest]) (*connect.Response[paymentv1.CreatePaymentResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	result, err := s.service.CreatePayment(ctx, payment.CreatePaymentParams{
		UserID:        claims.UserID,
		Address:       req.Msg.Address,
		PaymentMethod: convertPaymentMethod(req.Msg.Method),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CreatePaymentResponse{
		RequestId: req.Msg.RequestId,
		Url:       result.Url,
		Payment:   convertPaymentToProto(result.Payment),
	}), nil
}

// UpdatePayment implements the UpdatePayment method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) UpdatePayment(ctx context.Context, req *connect.Request[paymentv1.UpdatePaymentRequest]) (*connect.Response[paymentv1.UpdatePaymentResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	var method *model.PaymentMethod
	if req.Msg.Method != nil {
		method = util.ToPtr(convertPaymentMethod(*req.Msg.Method))
	}

	err = s.service.UpdatePayment(ctx, payment.UpdatePaymentParams{
		ID:      req.Msg.Id,
		UserID:  claims.UserID,
		Method:  method,
		Address: req.Msg.Address,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.UpdatePaymentResponse{}), nil
}

// CancelPayment implements the CancelPayment method from PaymentServiceHandler
func (s *ImplementedPaymentServiceHandler) CancelPayment(ctx context.Context, req *connect.Request[paymentv1.CancelPaymentRequest]) (*connect.Response[paymentv1.CancelPaymentResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.CancelPayment(ctx, payment.CancelPaymentParams{
		UserID:    claims.UserID,
		PaymentID: req.Msg.Id,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CancelPaymentResponse{}), nil
}

// Helper functions to convert between model and proto types
func convertPaymentToProto(p model.Payment) *paymentv1.Payment {
	products := make([]*paymentv1.ProductOnPayment, 0, len(p.Products))
	for _, pop := range p.Products {
		products = append(products, &paymentv1.ProductOnPayment{
			ItemQuantity: &common_pb.ItemQuantityInt64{
				ItemId:   pop.ItemID,
				Quantity: pop.Quantity,
			},
			SerialIds:  pop.SerialIDs,
			Price:      pop.Price,
			TotalPrice: pop.TotalPrice,
		})
	}

	return &paymentv1.Payment{
		Id:          p.ID,
		UserId:      p.UserID,
		Address:     p.Address,
		Method:      convertPaymentMethodToProto(p.Method),
		Total:       p.Total,
		Status:      convertStatusToProto(p.Status),
		DateCreated: p.DateCreated,
		Products:    products,
	}
}

func convertRefundToProto(r model.Refund) *paymentv1.Refund {
	return &paymentv1.Refund{
		Id:          r.ID,
		PaymentId:   r.PaymentID,
		Method:      convertRefundMethodToProto(r.Method),
		Status:      convertStatusToProto(r.Status),
		Reason:      r.Reason,
		Address:     r.Address,
		DateCreated: r.DateCreated,
		DateUpdated: r.DateUpdated,
		Resources:   r.Resources,
	}
}

func convertPaymentMethod(protoMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch protoMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CASH:
		return model.PaymentMethodCash
	case paymentv1.PaymentMethod_PAYMENT_METHOD_MOMO:
		return model.PaymentMethodMomo
	case paymentv1.PaymentMethod_PAYMENT_METHOD_VNPAY:
		return model.PaymentMethodVnpay
	case paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		panic("payment method is unspecified")
	default:
		panic("unknown payment method")
	}
}

func convertRefundMethod(protoMethod paymentv1.RefundMethod) model.RefundMethod {
	switch protoMethod {
	case paymentv1.RefundMethod_REFUND_METHOD_DROP_OFF:
		return model.RefundMethodDropOff
	case paymentv1.RefundMethod_REFUND_METHOD_PICK_UP:
		return model.RefundMethodPickUp
	case paymentv1.RefundMethod_REFUND_METHOD_UNSPECIFIED:
		panic("refund method is unspecified")
	default:
		panic("unknown refund method")
	}
}

func convertStatus(protoStatus common_pb.Status) model.Status {
	switch protoStatus {
	case common_pb.Status_STATUS_PENDING:
		return model.StatusPending
	case common_pb.Status_STATUS_SUCCESS:
		return model.StatusSuccess
	case common_pb.Status_STATUS_CANCELLED:
		return model.StatusCancelled
	case common_pb.Status_STATUS_FAILED:
		return model.StatusFailed
	case common_pb.Status_STATUS_UNSPECIFIED:
		panic("status is unspecified")
	default:
		panic("unknown status")
	}
}

func convertRefundMethodToProto(method model.RefundMethod) paymentv1.RefundMethod {
	switch method {
	case model.RefundMethodDropOff:
		return paymentv1.RefundMethod_REFUND_METHOD_DROP_OFF
	case model.RefundMethodPickUp:
		return paymentv1.RefundMethod_REFUND_METHOD_PICK_UP
	default:
		return paymentv1.RefundMethod_REFUND_METHOD_UNSPECIFIED
	}
}

func convertPaymentMethodToProto(method model.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case model.PaymentMethodCash:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CASH
	case model.PaymentMethodMomo:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_MOMO
	case model.PaymentMethodVnpay:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_VNPAY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func convertStatusToProto(status model.Status) common_pb.Status {
	switch status {
	case model.StatusPending:
		return common_pb.Status_STATUS_PENDING
	case model.StatusSuccess:
		return common_pb.Status_STATUS_SUCCESS
	case model.StatusCancelled:
		return common_pb.Status_STATUS_CANCELLED
	case model.StatusFailed:
		return common_pb.Status_STATUS_FAILED
	default:
		return common_pb.Status_STATUS_UNSPECIFIED
	}
}
