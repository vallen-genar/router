package router

import (
	"github.com/vallen-genar/router/entity"
	"net/http"
	"regexp"
)

type Route struct {
	path string
	middleware []entity.IController
	pathController entity.IController
	matchRegs map[*regexp.Regexp][]string
}

func NewRoute() *Route {
	return &Route{
		matchRegs: map[*regexp.Regexp][]string{},
	}
}

func (r *Route) AddMatch(match string, paramNames []string) {
	r.matchRegs[regexp.MustCompile(match)] = paramNames
}

func (r *Route) Middleware(controllers ...entity.IController) *Route {
	for _, c := range controllers {
		r.middleware = append(r.middleware, c)
	}
	return r
}

func (r *Route) Controller(controller entity.IController) *Route {
	r.pathController = controller
	return r
}

func (r *Route) Process(w http.ResponseWriter, req *http.Request, params map[string]string) {
	// process middleware
	for _, m := range r.middleware {
		err := m.Process(w, req, params)
		if err != nil {
			return
		}
	}

	// process controllers
	err := r.pathController.Process(w, req, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
