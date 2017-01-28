package api

import "github.com/humpback/humpback-center/ctrl"
import "github.com/gorilla/mux"

import (
	"net/http"
)

type handler func(c *Context) error

var routes = map[string]map[string]handler{
	"GET": {
		"/v1/_ping":                     ping,
		"/v1/cluster/groups":            getClusterGroups,
		"/v1/cluster/groups/{groupid}":  getClusterGroup,
		"/v1/repository/images/catalog": getRepositoryImagesCatalog,
		"/v1/repository/images/tags/*":  getRepositoryImagesTags,
	},
	"POST": {
		"/v1/cluster/groups/event":      postClusterGroupEvent,
		"/v1/repository/images/migrate": postRepositoryImagesMigrate,
	},
	"DELETE": {
		"/v1/repository/images/{name:.*}": deleteRepositoryImages,
	},
}

func NewRouter(controller *ctrl.Controller, enableCors bool) *mux.Router {

	router := mux.NewRouter()
	for method, mappings := range routes {
		for route, handler := range mappings {
			routemethod := method
			routepattern := route
			routehandler := handler
			wrap := func(w http.ResponseWriter, r *http.Request) {
				if enableCors {
					writeCorsHeaders(w, r)
				}
				c := NewContext(w, r, controller)
				routehandler(c)
			}
			router.Path(routepattern).Methods(routemethod).HandlerFunc(wrap)
			if enableCors {
				optionsmethod := "OPTIONS"
				optionshandler := optionsHandler
				wrap := func(w http.ResponseWriter, r *http.Request) {
					if enableCors {
						writeCorsHeaders(w, r)
					}
					c := NewContext(w, r, controller)
					optionshandler(c)
				}
				router.Path(routepattern).Methods(optionsmethod).HandlerFunc(wrap)
			}
		}
	}
	return router
}

func ping(ctx *Context) error {

	return ctx.JSON(http.StatusOK, "PANG")
}

func optionsHandler(ctx *Context) error {

	ctx.WriteHeader(http.StatusOK)
	return nil
}