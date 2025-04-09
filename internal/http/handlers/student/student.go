package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/goracijCerv/students-api/internal/storage"
	"github.com/goracijCerv/students-api/internal/types"
	"github.com/goracijCerv/students-api/internal/utils/response"
)

func Welcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la students api"))
	}
}

func New(storage storage.Storage) http.HandlerFunc {

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

		//Incercion de los datos
		id, err := storage.CreateStudent(student.Name, student.LastName, student.Email, student.Number, student.Age)
		if err != nil {
			slog.Error("user dont created something went wrong")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("user created suceessfuly", slog.String("userId", fmt.Sprint(id)))
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("obteniendo la información del usuario", slog.String("userId", fmt.Sprint(id)))

		num, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("id tiene un formato invalido")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		st, err := storage.GetStudentById(int64(num))
		if err != nil {
			slog.Error("Error al intentar obtener información")
			if err.Error() == "no student found" {
				response.WriteJson(w, http.StatusNotFound, err)
				return
			}

			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("Se ha optenido informacion del usuario correctamente", slog.String("userId", fmt.Sprint(id)))
		response.WriteJson(w, http.StatusFound, st)
	}
}
