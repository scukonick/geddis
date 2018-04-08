package geddis

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scukonick/geddis/db"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

type ServerAPI struct {
	store   *db.GeddisStore
	strings *stringsAPI
	arrays  *arraysAPI
	maps    *mapsAPI
	common  *commonAPI
}

func NewServerAPI(store *db.GeddisStore) *ServerAPI {
	return &ServerAPI{
		store:   store,
		strings: NewStringAPI(store),
		common:  NewCommonAPI(store),
		arrays:  newArraysAPI(store),
		maps:    newmapsAPI(store),
	}
}

func (api *ServerAPI) getRoutes() Routes {
	return Routes{
		Route{
			"Index",
			"GET",
			"/",
			Index,
		},

		// Strings
		Route{
			"GetStringByKey",
			"GET",
			"/strings/{key}",
			api.strings.GetStringByKey,
		},

		Route{
			"StringsKeyPost",
			"POST",
			"/strings/{key}",
			api.strings.StringsKeyPost,
		},

		// Arrays
		Route{
			"GetArrayByKey",
			"GET",
			"/arrays/{key}",
			api.arrays.GetByKey,
		},

		Route{
			"SetArray",
			"POST",
			"/arrays/{key}",
			api.arrays.Post,
		},

		Route{
			"GetArrayByKeyIndex",
			"GET",
			"/arrays/{key}/{index}",
			api.arrays.GetByKeyIndex,
		},

		// Maps
		Route{
			"GetMap",
			"GET",
			"/maps/{key}",
			api.maps.GetByKey,
		},

		Route{
			"SetMap",
			"POST",
			"/maps/{key}",
			api.maps.Post,
		},

		Route{
			"GetMapByKeySubkey",
			"GET",
			"/maps/{key}/{subkey}",
			api.maps.GetByKeySubKey,
		},

		// Common
		Route{
			"DeleteValue",
			"DELETE",
			"/common/{key}",
			api.common.DeleteValue,
		},
	}
}

func (api *ServerAPI) GetRouter() *mux.Router {
	routes := api.getRoutes()
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
