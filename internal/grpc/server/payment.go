package server

import (
	"context"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/grpc/server/interceptor"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service"
	"strconv"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServer
	paymentService *service.PaymentService
}

func NewPaymentServer(paymentService *service.PaymentService) *PaymentServer {
	return &PaymentServer{
		paymentService: paymentService,
	}
}

func (s *PaymentServer) Create(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	payment, err := s.paymentService.CreatePayment(ctx, service.CreatePaymentParams{
		UserID:        claims.UserID,
		Address:       req.Address,
		PaymentMethod: convertPaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId:  strconv.FormatInt(payment.Payment.ID, 10),
		PaymentUrl: payment.Url,
	}, nil
}

func convertPaymentMethod(protoMethod pb.PaymentMethod) model.PaymentMethod {
	switch protoMethod {
	case pb.PaymentMethod_CASH:
		return model.PaymentMethodCash
	case pb.PaymentMethod_MOMO:
		return model.PaymentMethodMomo
	case pb.PaymentMethod_VNPAY:
		return model.PaymentMethodVnpay
	case pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		panic("payment method is unspecified")
	default:
		panic("unknown payment method")
	}
}

func convertStatus(protoStatus pb.Status) model.Status {
	switch protoStatus {
	case pb.Status_PENDING:
		return model.StatusPending
	case pb.Status_SUCCESS:
		return model.StatusSuccess
	case pb.Status_CANCELLED:
		return model.StatusCancelled
	case pb.Status_FAILED:
		return model.StatusFailed
	case pb.Status_STATUS_UNSPECIFIED:
		panic("status is unspecified")
	default:
		panic("unknown status")
	}
}
