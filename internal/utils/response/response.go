package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK     = "ok"
	StastusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StastusError,
		Error:  err.Error(),
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	return Response{
		Status: StastusError,
		Error:  err.Error(),
	}
}
