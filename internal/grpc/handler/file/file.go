package file

import (
	"bytes"
	"context"
	"fmt"
	"shopnexus-go-service/internal/grpc/interceptor/auth"
	"shopnexus-go-service/internal/model"
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

// func (s *ImplementedFileServiceHandler) Upload(ctx context.Context, req *connect.ClientStream[filev1.UploadRequest]) (*connect.Response[filev1.UploadResponse], error) {
// 	claims, ok := ctx.Value(auth.CtxServerAccount).(model.Claims)
// 	if !ok {
// 		return nil, model.ErrTokenInvalid
// 	}

// 	var (
// 		filename string
// 	)

// 	// Initialize pipe for streaming data directly to S3
// 	pr, pw := io.Pipe()
// 	defer pr.Close()
// 	defer pw.Close()

// 	for req.Receive() {
// 		chunk := req.Msg()

// 		if filename == "" {
// 			filename = chunk.Name
// 		}

// 		// Write chunk data to file
// 		// Write file data to the pipe
// 		_, err := pw.Write(chunk.Content)
// 		if err != nil {
// 			pw.CloseWithError(err)
// 			return nil, fmt.Errorf("write to pipe: %v", err)
// 		}
// 	}

// 	// Check for errors in the stream
// 	if err := req.Err(); err != nil {
// 		return nil, fmt.Errorf("receive stream: %v", err)
// 	}

// 	// Upload to s3
// 	url, err := s.service.Upload(ctx, s3.GenKey(claims.UserID, filename), pr, false)
// 	if err != nil {
// 		return nil, fmt.Errorf("upload to s3: %v", err)
// 	}

// 	return connect.NewResponse(&filev1.UploadResponse{
// 		Url: url,
// 	}), nil
// }

func (s *ImplementedFileServiceHandler) Upload(ctx context.Context, req *connect.Request[filev1.UploadRequest]) (*connect.Response[filev1.UploadResponse], error) {
	claims, ok := ctx.Value(auth.CtxServerAccount).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
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
