package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goracijCerv/students-api/internal/types"
	"github.com/goracijCerv/students-api/internal/utils/response"
)

func Welcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la students api"))
	}
}

func New() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		slog.Info("creando un nuevo usuario")

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		//Validacion de la petición
		if err = validator.New().Struct(student); err != nil {
			validateErros := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErros))
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}

//https://youtu.be/OGhQhFKvMiM?t=5716
