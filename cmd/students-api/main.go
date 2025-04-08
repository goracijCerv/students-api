package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goracijCerv/students-api/internal/config"
	"github.com/goracijCerv/students-api/internal/http/handlers/student"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup a la base de datos
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", student.Welcome())
	router.HandleFunc("POST /api/students", student.New())
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
