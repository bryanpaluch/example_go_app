package controller

import (
	"encoding/json"
	"github.com/bryanpaluch/example_go_app/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"strconv"
)

type Router struct {
	router chi.Router
	db     db.DB
}

func NewRouter(d db.DB) (*Router, error) {
	return &Router{nil, d}, nil
}

func (r *Router) Start() {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/person/{id}", r.GetPersonByID)
	http.ListenAndServe(":8080", mux)
}

func (r *Router) GetPersonByID(w http.ResponseWriter, req *http.Request) {
	log.Println("get person by id hit")
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

func mustJSON(i interface{}) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return b
}
