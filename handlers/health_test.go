package handlers_test

import (
	"canvas/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestHealthz(t *testing.T) {
	t.Run("check body for /healthz", func(t *testing.T) {
		is := is.New(t)

		// init routes
		r := mux.NewRouter()
		r.HandleFunc("/healthz", handlers.Healthz) // test our health function

		_, _, body := makeRequest(r)
		is.Equal("up", body)
	})

	t.Run("returns 200 for /healthz", func(t *testing.T) {
		is := is.New(t)

		// init routes
		r := mux.NewRouter()
		r.HandleFunc("/healthz", handlers.Healthz) // test our health function

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