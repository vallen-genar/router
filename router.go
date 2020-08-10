package router

import (
	"github.com/g8y3e/router/controller"
	"github.com/g8y3e/router/entity"
	"net/http"
	"regexp"
)

const (
	defaultUriVarExp = "\\{(\\w+)\\}"
)

type Router struct {
	prefix string
	uriVarReg *regexp.Regexp

	getRoutes []*Route
	postRoutes []*Route
	putRoutes []*Route
	deleteRoutes []*Route

	httpNotFound entity.IController
}

func New(cf *Config) *Router {
	httpNotFound := cf.HttpNotFound
	if httpNotFound == nil {
		httpNotFound = &controller.HttpNotFound{}
	}

	uriVarExp := cf.UriVarExp
	if len(cf.UriVarExp) == 0{
		uriVarExp = defaultUriVarExp
	}

	return &Router{
		prefix: cf.Prefix,
		uriVarReg: regexp.MustCompile(uriVarExp),
		httpNotFound: httpNotFound,
		getRoutes:[]*Route{},
		postRoutes: []*Route{},
		putRoutes: []*Route{},
		deleteRoutes: []*Route{},
	}
}

func (r *Router) generatePathData(path string) (string, []string) {
	matchedParams := r.uriVarReg.FindAllString(path, -1)

	return r.uriVarReg.ReplaceAllString(path + "\\/*$", "(\\w+)"), matchedParams
}

func (r *Router) generateRoute(paths ...string) *Route {
	pathRoute := NewRoute()
	for _, path := range paths {
		path, pathParams := r.generatePathData(path)
		pathRoute.AddMatch(path, pathParams)
	}
	return pathRoute
}

func (r *Router) Get(paths ...string) *Route {
	pathRoute := r.generateRoute(paths...)
	r.getRoutes = append(r.getRoutes, pathRoute)
	return pathRoute
}

func (r *Router) Post(paths ...string) *Route {
	pathRoute := r.generateRoute(paths...)
	r.postRoutes = append(r.postRoutes, pathRoute)
	return pathRoute
}

func (r *Router) Put(paths ...string) *Route {
	pathRoute := r.generateRoute(paths...)
	r.putRoutes = append(r.putRoutes, pathRoute)
	return pathRoute
}

func (r *Router) Delete(paths ...string) *Route {
	pathRoute := r.generateRoute(paths...)
	r.deleteRoutes = append(r.deleteRoutes, pathRoute)
	return pathRoute
}

func (r *Router) Match(req *http.Request) (*Route, map[string]string) {
	// check routes type
	var searchRoutes []*Route
	if req.Method == http.MethodGet {
		searchRoutes = r.getRoutes
	} else if req.Method == http.MethodPost {
		searchRoutes = r.postRoutes
	} else if req.Method == http.MethodPut {
		searchRoutes = r.putRoutes
	} else if req.Method == http.MethodDelete {
		searchRoutes = r.deleteRoutes
	} else if searchRoutes == nil {
		return nil, nil
	}

	if len(req.URL.Path) < len(r.prefix) {
		return nil, nil
	}
	uri := req.URL.Path[len(r.prefix):]
	for _, searchRoute := range searchRoutes {
		// route /name/{variable}/{id}
		for matchReg, paramNames := range searchRoute.matchRegs {
			if (len(matchReg.String()) == 4 && len(uri) == 0) ||
				(len(matchReg.String()) != 4 && matchReg.MatchString(uri)) {
				subMatches := matchReg.FindAllStringSubmatch(uri, -1)
				params := map[string]string{}
				if len(subMatches) > 0 && len(subMatches[0]) > len(paramNames) {
					for index, value := range paramNames {
						params[value[1:len(value)-1]] = subMatches[0][index + 1]
					}
				}

				return searchRoute, params
			}
		}
	}
	return nil, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	routeMatch, params := r.Match(req)
	if routeMatch == nil {
		r.httpNotFound.Process(w, req, params)
	} else {
		routeMatch.Process(w, req, params)
	}
}
