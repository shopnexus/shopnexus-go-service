package common

import (
	"shopnexus-go-service/internal/model"

	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/common"
)

func ToProtoPaginationResponse[T any](data model.PaginateResult[T]) *common.PaginationResponse {
	return &common.PaginationResponse{
		Page:       data.Page,
		Limit:      data.Limit,
		Total:      uint32(data.Total),
		NextPage:   data.NextPage,
		NextCursor: data.NextCursor,
	}
}
