package entity

import "net/http"

type IController interface {
	 Process(w http.ResponseWriter, req *http.Request, params map[string]string) error
}
