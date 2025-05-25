package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/model"
	"strconv"

	"go.uber.org/zap"
)

func writeError(w http.ResponseWriter, errCode string, message string) {
	response, err := json.Marshal(CommonResponse{
		Data: nil,
		Error: &Error{
			Code:    errCode,
			Message: message,
		},
	})
	if err != nil {
		w.Write([]byte("Error marshalling JSON"))
		return
	}
	w.Write(response)
}

func writeResponse(w http.ResponseWriter, dto any) {
	response, err := json.Marshal(dto)
	if err != nil {
		writeError(w, http.StatusText(http.StatusInternalServerError), "Error marshalling JSON")
		return
	}

	w.Write(response)
}

func FromDTO(w http.ResponseWriter, dto any, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	writeResponse(w, CommonResponse{
		Data:  dto,
		Error: nil,
	})
}

func FromError(w http.ResponseWriter, err error, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	// Internal server error
	var errWithCode *model.ErrorWithCode
	if errors.As(err, &errWithCode) {
		logger.Log.Error("error", zap.Error(errWithCode.Err))
		writeError(w, errWithCode.Code, errWithCode.Msg)
		return
	}

	// Normal http error
	logger.Log.Error("error", zap.Error(err))
	writeError(w, strconv.Itoa(httpCode), err.Error())
}

func FromHTTPError(w http.ResponseWriter, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	writeError(w, strconv.Itoa(httpCode), http.StatusText(httpCode))
}

func FromPaginate[T any](w http.ResponseWriter, paginateResult model.PaginateResult[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	writeResponse(w, PaginateResponse[T]{
		Data: paginateResult.Data,
		Pagination: Pagination{
			Limit:      paginateResult.Limit,
			Page:       paginateResult.Page,
			Total:      paginateResult.Total,
			NextPage:   paginateResult.NextPage,
			NextCursor: paginateResult.NextCursor,
		},
		Error: nil,
	})
}
