package binding

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var (
	validate = validator.New()
	decoder  = schema.NewDecoder()
)

// For url.Values (map[string][]string); `schema` and `validate` tags are used
func BindQuery(dst any, src url.Values) error {
	if err := decoder.Decode(dst, src); err != nil {
		return err
	}
	return validate.Struct(dst)
}

// For JSON body (io.Reader); `json` and `validate` tags are used
func BindJSON(dst any, src io.Reader) error {
	if err := json.NewDecoder(src).Decode(dst); err != nil {
		return err
	}
	return validate.Struct(dst)
}
