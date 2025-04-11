// @title Students API
// @version 1.0
// @description Bueno esta es una api, para practicar buenas practicas en go y que entoeria esta hecha de forma que es escalable.
// @host localhost:8082
// @BasePath /api/
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/goracijCerv/students-api/docs"
	"github.com/goracijCerv/students-api/internal/config"
	"github.com/goracijCerv/students-api/internal/http/handlers/student"
	"github.com/goracijCerv/students-api/internal/storage/sqlite"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup a la base de datos
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Se ha inicializado corractamente la bd", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", student.Welcome())
	router.HandleFunc("POST /api/student", student.New(storage))
	router.HandleFunc("GET /api/student/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/student", student.GetListStudents(storage))
	router.HandleFunc("PUT /api/student/{id}", student.UpdateById(storage))
	router.HandleFunc("DELETE /api/student/{id}", student.DeleteStudent(storage))
	// Serve Swagger UI
	swaggerURL := fmt.Sprintf("http://%s/swagger/doc.json", cfg.Addr) // Construct the URL to swagger.json dynamically
	router.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL(swaggerURL),
	))
	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Servidor iniciado en el puerto ", slog.String("adress", cfg.Addr))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("fail to start server")
		}
	}()
	<-done

	slog.Info("apagagando el servidor")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Hubo un problema al tratar de cerar el servidor ", slog.String("error", err.Error()))
	}

	slog.Info("apagado del server correcto")

}

//Para pode correr el programa go run cmd/students-api/main.go -config config/local.yaml

//Paquetes que se utilizaron para la creacion de la documentacion de swagger
// swag init -g cmd/students-api/main.go genera la documentacion de swagger
// go get -u github.com/swaggo/http-swagger
// go install github.com/swaggo/swag/cmd/swag@latest
