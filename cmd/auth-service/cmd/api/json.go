package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_578 //1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if r.Body == nil {
		return errors.New("request body is nil")
	}
	err := decoder.Decode(data)
	if err != nil {
		if err == io.EOF {
			return errors.New("request body is empty")

		}
		return err
	}

	return nil

}

func writeJsonError(w http.ResponseWriter, status int, message string, errors []string) error {
	type envelope struct {
		Errors  []string `json:"errors"`
		Message string   `json:"message"`
	}
	return writeJson(w, status, &envelope{Message: message, Errors: errors})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, message string, data any) error {
	type envelope struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	return writeJson(w, status, &envelope{Message: message, Data: data})
}
