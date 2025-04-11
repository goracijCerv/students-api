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

// Welcome godoc
// @Summary De vuelve un mensaje de bienvenida
// @Description Da un mensaje de bienvenida a la api
// @Tags welcome
// @Success 200 {string} string "Bienvenido a la students api"
// @Router / [get]
func Welcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la students api"))
	}
}

// CreateStudent godoc
// @Summary Crea un nuevo estudiante
// @Description Obtiene los datos de la estructura student y crea un nuevo estudiante
// @Tags student
// @Accept json
// @Produce json
// @Param student body types.Student true "Un objeto de tipo estudiante o student y este nos servira para crear no es necesario el id"
// @Success 201 {object} map[string]string "Mensjae de que todo salio bien"
// @Failure 400 {object} response.Response "Cuerpo de  la petición no valido"
// @Failure 500 {object} response.Response "Error del servidor"
// @Router /api/student [post]
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

// getByid godoc
// @Summary Obtiene la informacion de un estudiante por id
// @Description Obtiene los datos del estudiante espesificado por id
// @Tags student
// @Accept json
// @Produce json
// @Param id path int true "id del estudiante"
// @Success 200 {object} types.Student
// @Failure 400 {object} response.Response "Error en el path"
// @Failure 404 {object} response.Response "No se encontro el recurso"
// @Failure 500 {object} response.Response "Error del servidor"
// @Router /api/student/{id} [get]
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

			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Se ha optenido informacion del usuario correctamente", slog.String("userId", fmt.Sprint(id)))
		response.WriteJson(w, http.StatusOK, st)
	}
}

// GetAllStudents godoc
// @Summary Obtiene todos los estudiantes
// @Description Obtiene todos los datos de los estudiantes
// @Tags student
// @Accept json
// @Produce json
// @Success 200 {array} types.Student "Lista de estudiantes"
// @Failure 500 {object} response.Response "Error del servidor"
// @Router /api/student [get]
func GetListStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Obteniendo la lista de estudiantes")
		students, err := storage.GetAllStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}

		slog.Info("Se han obtenido todos los usuarios correctamente")
		response.WriteJson(w, http.StatusOK, students)
	}
}

// UpdateStudent godoc
// @Summary Le hace un update al student
// @Description se actualiza la informacioón del cliente
// @Tags student
// @Accept json
// @Produce json
// @Param id path int true "id del estudiante"
// @Param  student body types.Student true "la nueva informacion del student"
// @Success 200 {object} map[string]string "Mensaje de que salio bien"
// @Failure 400 {object} response.Response "Error en el path"
// @Failure 500 {object} response.Response "Error del servidor"
// @Router /api/student/{id} [put]
func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("editando el usario con el id", slog.String("userId", fmt.Sprint(id)))

		num, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("id tiene un formato invalido")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		var student types.Student
		err = json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			slog.Error("error en json decoder")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		student.Id = int64(num)
		//Validacion de la petición
		if err = validator.New().Struct(student); err != nil {
			slog.Error("error de validación")
			validateErros := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErros))
			return
		}

		//Updateo de  usuario
		err = storage.UpdateStudent(int64(num), student.Name, student.LastName, student.Email, student.Number, student.Age)
		if err != nil {
			slog.Error("Error en el querry")
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("La información se updateo correctamente")
		response.WriteJson(w, http.StatusOK, student)

	}
}

// DeleteStudent godoc
// @Summary Borra un estudiante
// @Description Borra un estudiante con el id que se le pasa
// @Tags student
// @Accept json
// @Produce json
// @Param id path int true "id del estudiante"
// @Success 200 {object} map[string]string "Mensaje de que salio bien"
// @Failure 400 {object} response.Response "Error en el path"
// @Failure 500 {object} response.Response "Error del servidor"
// @Router /api/student/{id} [DELETE]
func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Borrando un elemento ", slog.String("userId", fmt.Sprint(id)))

		num, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("id no tiene un formato valido")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = storage.DeleteStudentById(int64(num))
		if err != nil {
			slog.Error("Error in the query")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Se borro correctamente")
		response.WriteJson(w, http.StatusOK, map[string]string{"succes": "OK"})
	}
}

// https://youtu.be/OGhQhFKvMiM?t=9475
