package controller

import (
	"encoding/json"
	"github.com/bryanpaluch/example_go_app/example"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Router struct {
	chi.Router
	db example.DB
}

func NewRouter(d example.DB) (*Router, error) {
	mux := chi.NewRouter()
	r := &Router{mux, d}

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/person/{id}", r.GetPersonByID)
	mux.Post("/person", r.InsertNewPerson)

	return r, nil
}

func (r *Router) Start() {
	http.ListenAndServe(":8080", r)
}

func (r *Router) GetPersonByID(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	p, err := r.db.GetPersonByID(req.Context(), idNum)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if p == nil {
		w.WriteHeader(404)
		return
	}

	w.Write(mustJSON(p))
}

func (r *Router) InsertNewPerson(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var person example.Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = r.db.AddPerson(req.Context(), &person)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(mustJSON(person))
}

func mustJSON(i interface{}) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return b
}
