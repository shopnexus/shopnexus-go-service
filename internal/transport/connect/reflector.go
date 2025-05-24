package connect

import (
	"net/http"

	"connectrpc.com/grpcreflect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/file/v1/filev1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1/paymentv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1/productv1connect"
)

func InitReflector(mux *http.ServeMux) {
	reflector := grpcreflect.NewStaticReflector(
		filev1connect.FileServiceName,
		accountv1connect.AccountServiceName,
		productv1connect.ProductServiceName,
		paymentv1connect.PaymentServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}
