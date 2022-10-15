package handlers

import (
	"fmt"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "up")
}
