package geddis

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scukonick/geddis/db"
	"github.com/scukonick/geddis/serverxxx"
)

type commonAPI struct {
	store *db.GeddisStore
}

func newCommonAPI(s *db.GeddisStore) *commonAPI {
	return &commonAPI{
		store: s,
	}
}

func (s *commonAPI) deleteValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.store.Del(key)

	w.WriteHeader(http.StatusOK)
}

func (s *commonAPI) keys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if key == "*" {
		key = ""
	}

	resp := &swagger.Array{
		Values: s.store.Keys(key),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Failed to sent result: %v", err)
	}
}
