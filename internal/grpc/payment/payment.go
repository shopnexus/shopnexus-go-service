package payment

import (
	"context"
	"shopnexus-go-service/internal/grpc/interceptor"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/payment"
	"strconv"

	"connectrpc.com/connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/common"
	paymentv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1/paymentv1connect"
)

var _ paymentv1connect.PaymentServiceHandler = (*PaymentServer)(nil)

type PaymentServer struct {
	paymentv1connect.UnimplementedPaymentServiceHandler
	service *payment.PaymentService
}

func NewPaymentServer(paymentService *payment.PaymentService) *PaymentServer {
	return &PaymentServer{
		service: paymentService,
	}
}

// GetPayment implements the GetPayment method from PaymentServiceHandler
func (s *PaymentServer) GetPayment(ctx context.Context, req *connect.Request[paymentv1.GetPaymentRequest]) (*connect.Response[paymentv1.GetPaymentResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	paymentID, err := strconv.ParseInt(req.Msg.PaymentId, 10, 64)
	if err != nil {
		return nil, err
	}

	payment, err := s.service.GetPayment(ctx, paymentID, claims.UserID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.GetPaymentResponse{
		Payment: convertPaymentToProto(payment),
	}), nil
}

// ListPayments implements the ListPayments method from PaymentServiceHandler
func (s *PaymentServer) ListPayments(ctx context.Context, req *connect.Request[paymentv1.ListPaymentsRequest]) (*connect.Response[paymentv1.ListPaymentsResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	payments, err := s.service.ListPayments(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	protoPayments := make([]*paymentv1.Payment, 0, len(payments))
	for _, p := range payments {
		protoPayments = append(protoPayments, convertPaymentToProto(p))
	}

	return connect.NewResponse(&paymentv1.ListPaymentsResponse{
		Payments: protoPayments,
	}), nil
}

// CreatePayment implements the CreatePayment method from PaymentServiceHandler
func (s *PaymentServer) CreatePayment(ctx context.Context, req *connect.Request[paymentv1.CreatePaymentRequest]) (*connect.Response[paymentv1.CreatePaymentResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	result, err := s.service.CreatePayment(ctx, payment.CreatePaymentParams{
		UserID:        claims.UserID,
		Address:       req.Msg.Address,
		PaymentMethod: convertPaymentMethod(req.Msg.PaymentMethod),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CreatePaymentResponse{
		PaymentId:  strconv.FormatInt(result.Payment.ID, 10),
		PaymentUrl: result.Url,
	}), nil
}

// UpdatePayment implements the UpdatePayment method from PaymentServiceHandler
func (s *PaymentServer) UpdatePayment(ctx context.Context, req *connect.Request[paymentv1.UpdatePaymentRequest]) (*connect.Response[paymentv1.UpdatePaymentResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	paymentID, err := strconv.ParseInt(req.Msg.PaymentId, 10, 64)
	if err != nil {
		return nil, err
	}

	payment, err := s.service.UpdatePayment(ctx, payment.UpdatePaymentParams{
		ID:     paymentID,
		UserID: claims.UserID,
		Status: convertStatus(req.Msg.Status),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.UpdatePaymentResponse{
		Payment: convertPaymentToProto(payment),
	}), nil
}

// CancelPayment implements the CancelPayment method from PaymentServiceHandler
func (s *PaymentServer) CancelPayment(ctx context.Context, req *connect.Request[paymentv1.CancelPaymentRequest]) (*connect.Response[paymentv1.CancelPaymentResponse], error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	paymentID, err := strconv.ParseInt(req.Msg.PaymentId, 10, 64)
	if err != nil {
		return nil, err
	}

	err = s.service.CancelPayment(ctx, paymentID, claims.UserID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&paymentv1.CancelPaymentResponse{}), nil
}

// Helper functions to convert between model and proto types
func convertPaymentToProto(p model.Payment) *paymentv1.Payment {
	products := make([]*paymentv1.ProductOnPayment, 0, len(p.Products))
	for _, product := range p.Products {
		products = append(products, &paymentv1.ProductOnPayment{
			Id:         product.ID,
			Quantity:   int32(product.Quantity),
			Price:      product.Price,
			TotalPrice: product.TotalPrice,
		})
	}

	return &paymentv1.Payment{
		Id:            strconv.FormatInt(p.ID, 10),
		UserId:        strconv.FormatInt(p.UserID, 10),
		Address:       p.Address,
		PaymentMethod: convertPaymentMethodToProto(p.PaymentMethod),
		Total:         p.Total,
		Status:        convertStatusToProto(p.Status),
		DateCreated:   p.DateCreated,
		Products:      products,
	}
}

func convertRefundToProto(r model.Refund) *paymentv1.Refund {
	var address string
	if r.Address != nil {
		address = *r.Address
	}

	return &paymentv1.Refund{
		Id:          strconv.FormatInt(r.ID, 10),
		PaymentId:   strconv.FormatInt(r.PaymentID, 10),
		Method:      convertRefundMethodToProto(r.Method),
		Status:      convertStatusToProto(r.Status),
		Reason:      r.Reason,
		Address:     address,
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

func convertStatus(protoStatus common.Status) model.Status {
	switch protoStatus {
	case common.Status_STATUS_PENDING:
		return model.StatusPending
	case common.Status_STATUS_SUCCESS:
		return model.StatusSuccess
	case common.Status_STATUS_CANCELLED:
		return model.StatusCancelled
	case common.Status_STATUS_FAILED:
		return model.StatusFailed
	case common.Status_STATUS_UNSPECIFIED:
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

func convertStatusToProto(status model.Status) common.Status {
	switch status {
	case model.StatusPending:
		return common.Status_STATUS_PENDING
	case model.StatusSuccess:
		return common.Status_STATUS_SUCCESS
	case model.StatusCancelled:
		return common.Status_STATUS_CANCELLED
	case model.StatusFailed:
		return common.Status_STATUS_FAILED
	default:
		return common.Status_STATUS_UNSPECIFIED
	}
}
