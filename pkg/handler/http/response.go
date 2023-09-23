package http

import (
	"encoding/json"
	"errors"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/apperrors"
	"net/http"
)

// Response in order to unify the returned response structure
type Response struct {
	Code   int  `json:"-"`
	Pretty bool `json:"-"`
	Data   any  `json:"data,omitempty"`
	Error  any  `json:"error,omitempty"`
}

func (a Response) JSON(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")

	if err, ok := a.Error.(error); ok {
		a.Error = err.Error()

		switch {
		case errors.Is(err, apperrors.ErrObjectNotFound):
			a.Code = http.StatusNotFound
		default:
			a.Code = http.StatusInternalServerError
		}
	}

	w.WriteHeader(a.Code)
	_ = json.NewEncoder(w).Encode(a)
}
