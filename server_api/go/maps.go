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

type mapsAPI struct {
	store *db.GeddisStore
}

func newmapsAPI(s *db.GeddisStore) *mapsAPI {
	return &mapsAPI{
		store: s,
	}
}

func (s *mapsAPI) GetByKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := s.store.GetMap(key)
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

	data := &swagger.MapValue{
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

func (s *mapsAPI) GetByKeySubKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subKey, ok := vars["subkey"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := s.store.GetByKey(key, subKey)
	switch {
	case err == db.ErrNotFound || err == db.ErrInvalidIndex:
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

func (s *mapsAPI) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

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

	body := &swagger.SetMapReq{}
	err = json.Unmarshal(input, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ttl := time.Duration(body.Ttl)

	s.store.SetMap(key, body.Value, ttl*time.Second)

	w.WriteHeader(http.StatusOK)
}
