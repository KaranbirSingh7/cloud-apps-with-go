package handlers_test

import (
	"canvas/handlers"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/matryer/is"
)

type pingerMock struct {
	err error
}

func (p *pingerMock) Ping(ctx context.Context) error {
	return p.err
}

func TestHealthz(t *testing.T) {
	t.Run("return 502 if DB is offline", func(t *testing.T) {
		is := is.New(t)

		// init routes
		r := mux.NewRouter()
		handlers.HealthzWrapper(r, &pingerMock{err: errors.New("DB offline")})

		code, _, _ := makeRequest(r)
		is.Equal(http.StatusBadGateway, code)

	})

	t.Run("returns 200 for /healthz", func(t *testing.T) {
		is := is.New(t)

		// init routes
		r := mux.NewRouter()
		handlers.HealthzWrapper(r, &pingerMock{})
		// r.HandleFunc("/healthz", handlers.Healthz) // test our health function

		statusCode, _, _ := makeRequest(r)
		is.Equal(http.StatusOK, statusCode)
	})
}

func makeRequest(r *mux.Router) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder() // mock http

	r.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return 0, nil, "" //aka error/failure decoding reading body
	}
	return result.StatusCode, result.Header, string(bodyBytes)

}
