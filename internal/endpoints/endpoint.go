package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"
	"orbit-app/pkg"
)

type Endpoint func(w http.ResponseWriter, r *http.Request) error

func (f Endpoint) Unwrap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if errors.As(err, &EndpointError{}) {
				endpointError := err.(EndpointError)
				w.WriteHeader(endpointError.StatusCode)
				w.Write([]byte(endpointError.Err))
				return
			} else if errors.As(err, &pkg.ValidationError{}) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(err.Error()))
		}
	}
}

type EndpointError struct {
	StatusCode int    `json:"status_code"`
	Err        string `json:"error"`
}

func (e EndpointError) Error() string {
	res, _ := json.Marshal(e)
	return string(res)
}
