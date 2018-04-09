package geddis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
	"github.com/scukonick/geddis/db"
	"github.com/scukonick/geddis/serverxxx"
)

type stringsAPI struct {
	store *db.GeddisStore
}

func newStringAPI(s *db.GeddisStore) *stringsAPI {
	return &stringsAPI{
		store: s,
	}
}

func (s *stringsAPI) getString(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := s.store.GetStr(key)
	switch {
	case err == db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	case err == db.ErrInvalidType:
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := &swagger.StringValue{
		Value: value,
	}
	js, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Printf("Failed to sent result: %v", err)
	}

}

func (s *stringsAPI) post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := &swagger.SetStringValueReq{}
	err = json.Unmarshal(input, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ttl := time.Duration(body.Ttl)
	s.store.SetStr(key, body.Value, ttl*time.Second)
}
