package geddis

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/scukonick/geddis/cli/lib"
	"github.com/scukonick/geddis/db"
)

type arraysAPI struct {
	store *db.GeddisStore
}

func newArraysAPI(s *db.GeddisStore) *arraysAPI {
	return &arraysAPI{
		store: s,
	}
}

func (s *arraysAPI) GetByKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := s.store.GetArr(key)
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

	data := &swagger.Array{
		Values: value,
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

func (s *arraysAPI) GetByKeyIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	indexStr, ok := vars["index"]
	if !ok {
		// should not happen as we have routing
		// but anyway it's better to check
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	value, err := s.store.GetByIndex(key, index)
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

func (s *arraysAPI) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

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

	body := &swagger.SetArrayReq{}
	err = json.Unmarshal(input, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ttl := time.Duration(body.Ttl)

	s.store.SetArr(key, body.Values, ttl*time.Second)

	w.WriteHeader(http.StatusOK)
}
