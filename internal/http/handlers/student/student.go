package student

import "net/http"

func Welcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la students api"))
	}
}

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la students api"))
	}
}
