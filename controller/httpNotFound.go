package controller

import (
	"fmt"
	"net/http"
)

type HttpNotFound struct {
}

func (h *HttpNotFound) Process(w http.ResponseWriter, req *http.Request, params map[string]string) error {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, `{"error": "Page not found"}`)
	return nil
}
