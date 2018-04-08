package geddis

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scukonick/geddis/db"
)

type commonAPI struct {
	store *db.GeddisStore
}

func NewCommonAPI(s *db.GeddisStore) *commonAPI {
	return &commonAPI{
		store: s,
	}
}

func (s *commonAPI) DeleteValue(w http.ResponseWriter, r *http.Request) {
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
