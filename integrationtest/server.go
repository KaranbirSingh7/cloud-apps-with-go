package integrationtest

import (
	"canvas/server"
	"net/http"
	"testing"
	"time"
)

// CreateServer for testing on port 8081, returning a cleanup function that stops the server.
// Usage:
//
//	cleanup := CreateServer()
//	defer cleanup()
func CreateServer() func() {
	db, cleanupDB := CreateDatabase()
	s := server.NewServer(server.Options{
		Host:     "localhost",
		Port:     8081,
		Database: db,
	})

	// run this in background and panic if there is any error
	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()

	// try to make calls to out server
	for {
		_, err := http.Get("http://localhost:8081/")
		if err == nil {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	// return function that can be called to stop server
	return func() {
		if err := s.Stop(); err != nil {
			panic(err)
		}
		cleanupDB()
	}
}

func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
}
