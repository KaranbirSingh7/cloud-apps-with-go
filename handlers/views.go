package handlers

import (
	"canvas/views"
	"net/http"
)


func IndexPage(w http.ResponseWriter, r *http.Request){
	_ = views.FrontPage().Render(w)
}
