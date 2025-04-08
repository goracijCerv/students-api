package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goracijCerv/students-api/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup a la base de datos
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la students api"))
	})
	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Println("Servidor iniciado en el puerto ", cfg.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("fail to strat server")
	}

}

//Para pode correr el programa go run cmd/students-api/main.go -config config/local.yaml
