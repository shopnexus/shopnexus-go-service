package file

import (
	"bytes"
	"context"
	"fmt"
	"shopnexus-go-service/internal/server/connect/interceptor/auth"
	"shopnexus-go-service/internal/service/s3"

	"connectrpc.com/connect"
	filev1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/file/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/file/v1/filev1connect"
)

type ImplementedFileServiceHandler struct {
	filev1connect.UnimplementedFileServiceHandler
	service *s3.S3Service
}

func NewFileServiceHandler(s3Service *s3.S3Service) filev1connect.FileServiceHandler {
	return &ImplementedFileServiceHandler{
		service: s3Service,
	}
}

func (s *ImplementedFileServiceHandler) Upload(ctx context.Context, req *connect.Request[filev1.UploadRequest]) (*connect.Response[filev1.UploadResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(req.Msg.Content)

	// Upload to s3
	url, err := s.service.Upload(ctx, s3.GenKey(claims.UserID, req.Msg.Name), body, false)
	if err != nil {
		return nil, fmt.Errorf("upload to s3: %v", err)
	}

	return connect.NewResponse(&filev1.UploadResponse{
		Url: url,
	}), nil
}
