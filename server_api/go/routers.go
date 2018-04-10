package geddis

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scukonick/geddis/db"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

// ServerAPI represents HTTP REST API for geddis
type ServerAPI struct {
	store   *db.GeddisStore
	strings *stringsAPI
	arrays  *arraysAPI
	maps    *mapsAPI
	common  *commonAPI
}

// NewServerAPI returns newly initialized ServerAPI instance
func NewServerAPI(store *db.GeddisStore) *ServerAPI {
	return &ServerAPI{
		store:   store,
		strings: newStringAPI(store),
		common:  newCommonAPI(store),
		arrays:  newArraysAPI(store),
		maps:    newmapsAPI(store),
	}
}

func (api *ServerAPI) getRoutes() routes {
	return routes{
		route{
			"index",
			"GET",
			"/",
			index,
		},

		// Strings
		route{
			"getStringByKey",
			"GET",
			"/strings/{key}",
			api.strings.getString,
		},

		route{
			"StringsKeyPost",
			"POST",
			"/strings/{key}",
			api.strings.post,
		},

		// Arrays
		route{
			"GetArrayByKey",
			"GET",
			"/arrays/{key}",
			api.arrays.GetByKey,
		},

		route{
			"SetArray",
			"POST",
			"/arrays/{key}",
			api.arrays.Post,
		},

		route{
			"GetArrayByKeyIndex",
			"GET",
			"/arrays/{key}/{index}",
			api.arrays.GetByKeyIndex,
		},

		// Maps
		route{
			"GetMap",
			"GET",
			"/maps/{key}",
			api.maps.GetByKey,
		},

		route{
			"SetMap",
			"POST",
			"/maps/{key}",
			api.maps.Post,
		},

		route{
			"GetMapByKeySubkey",
			"GET",
			"/maps/{key}/{subkey}",
			api.maps.GetByKeySubKey,
		},

		// Common
		route{
			"deleteValue",
			"DELETE",
			"/delete/{key}",
			api.common.deleteValue,
		},

		route{
			"GetKeys",
			"GET",
			"/keys/{key}",
			api.common.keys,
		},
	}
}

// GetRouter returns router from ServerAPI, so the client
// is able to call listen and stop as needed
func (api *ServerAPI) GetRouter() *mux.Router {
	routes := api.getRoutes()
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
