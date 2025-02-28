package response

import (
	"net/http"
	"shopnexus-go-service/internal/model"

	"github.com/bytedance/sonic"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CommonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   *Error `json:"error"`
}

type Pagination struct {
	Limit      int32   `json:"limit"` // Number of items per page
	Page       int32   `json:"page"`  // Current page
	Total      int64   `json:"total"` // Total number of items
	NextPage   *int32  `json:"next_page"`
	NextCursor *string `json:"next_cursor"`
}

type PaginateResponse struct {
	CommonResponse
	Pagination Pagination `json:"pagination"`
}

// FromError returns an error response with status code 500
func FromError(w http.ResponseWriter, err error) error {
	msg := "Internal Server Error"

	if err != nil {
		msg = err.Error()
	}

	w.WriteHeader(http.StatusInternalServerError)
	return sonic.ConfigFastest.NewEncoder(w).Encode(&CommonResponse{
		Status:  "error",
		Message: msg,
		Error: &Error{
			Code:    http.StatusInternalServerError,
			Message: msg,
		},
	})
}

// FromMessage returns a response with a message, success if code is 2xx, error otherwise
func FromMessage(w http.ResponseWriter, code int, message string) error {
	w.WriteHeader(code)

	if code >= 400 && code < 500 {
		return sonic.ConfigFastest.NewEncoder(w).Encode(&CommonResponse{
			Status:  "error",
			Message: message,
			Error: &Error{
				Code:    code,
				Message: message,
			},
		})
	} else {
		return sonic.ConfigFastest.NewEncoder(w).Encode(&CommonResponse{
			Status:  "success",
			Message: message,
		})
	}
}

func FromDTO(w http.ResponseWriter, code int, dto any) error {
	w.WriteHeader(code)

	if code >= 400 && code < 500 {
		return sonic.ConfigFastest.NewEncoder(w).Encode(&CommonResponse{
			Status:  "error",
			Message: "Error",
			Data:    dto,
		})
	} else {
		return sonic.ConfigFastest.NewEncoder(w).Encode(&CommonResponse{
			Status:  "success",
			Message: "Success",
			Data:    dto,
		})
	}
}

func FromPagination[T any](w http.ResponseWriter, code int, dto model.PaginateResult[T]) error {
	sonic.ConfigFastest.NewEncoder(w).Encode(&PaginateResponse{
		CommonResponse: CommonResponse{
			Status:  "success",
			Message: "Success",
			Data:    dto.Data,
		},
		Pagination: Pagination{
			Total:      dto.Total,
			Page:       dto.Page,
			Limit:      dto.Limit,
			NextPage:   dto.NextPage,
			NextCursor: dto.NextCursor,
		},
	})
	return nil
}
