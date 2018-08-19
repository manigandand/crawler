package main

import (
	"net/http"
	"strings"
)

func spiderman(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	websiteAddress := strings.TrimSpace(r.PostFormValue("website_address"))

	w.Write([]byte(websiteAddress))
}
