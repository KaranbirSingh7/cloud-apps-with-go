package model_test

import (
	"canvas/model"
	"testing"

	"github.com/matryer/is"
)

func TestEmail_IsValid(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"me@example.com", true},
		{"@example.com", false},
		{"me@", false},
		{"", false},
		{"@", false},
	}

	t.Run("reports valid email address", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.address, func(t *testing.T) {
				is := is.New(t)
				e := model.Email(test.address)
				is.Equal(test.valid, e.IsValid())
			})
		}
	})
	// this is just a fun test to check if total test coverage goes high
	// if you are wondering, it does increase total test converage
	t.Run("check if email string matches with user input", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.address, func(t *testing.T) {
				is := is.New(t)
				e := model.Email(test.address)
				is.Equal(e.String(), test.address)
			})
		}
	})
}
