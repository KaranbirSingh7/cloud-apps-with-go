package handlers_test

import (
	"canvas/handlers"
	"canvas/model"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/matryer/is"
)

// following is mock struct
type signupperMock struct {
	email model.Email
}

// this function would be called instead of our other 'SignupForNewsletter' because testing and mocking
func (s *signupperMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	s.email = email //add something here to know when this func is being called during test
	return "", nil
}

// TODO: add TDD (table driven tests) here
func TestNewsletterSignup(t *testing.T) {
	r := mux.NewRouter()
	s := &signupperMock{}

	// mux router and signup interface
	handlers.NewsletterSignupWrapper(r, s)

	t.Run("signs up a valid email address", func(t *testing.T) {
		is := is.New(t)
		code, _, _ := makePostRequest(r, "/newsletter/signup", createFormHeader(), strings.NewReader(
			"email=me%40example.com"))
		is.Equal(http.StatusFound, code)
		is.Equal(model.Email("me@example.com"), s.email)
	})

	t.Run("rejects an invalid email address", func(t *testing.T) {
		is := is.New(t)
		code, _, _ := makePostRequest(r, "/newsletter/signup", createFormHeader(), strings.NewReader(
			"email=measdads.com"))
		is.Equal(http.StatusBadRequest, code)
	})

}

func makePostRequest(handler http.Handler, target string, header http.Header, body io.Reader) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header = header
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	return result.StatusCode, result.Header, string(bodyBytes)
}

func createFormHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	return header
}
