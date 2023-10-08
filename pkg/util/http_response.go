package util

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/apperrors"
)

// Response in order to unify the returned response structure
type Response struct {
	Code   int    `json:"-"`
	Pretty bool   `json:"-"`
	Data   any    `json:"data,omitempty"`
	Err    string `json:"error,omitempty"`
}

func (a Response) Error() string {
	return a.Err
}

func (a Response) JSON(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")

	if a.Err != "" {
		switch {
		case errors.Is(a, apperrors.ErrObjectNotFound):
			a.Code = http.StatusNotFound
		default:
			a.Code = http.StatusInternalServerError
		}
	}

	w.WriteHeader(a.Code)
	_ = json.NewEncoder(w).Encode(a)
}
