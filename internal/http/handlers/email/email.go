package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goracijCerv/students-api/internal/storage"
	"github.com/goracijCerv/students-api/internal/types"
	"github.com/goracijCerv/students-api/internal/utils"
	"github.com/goracijCerv/students-api/internal/utils/response"
)

// sendEmail godoc
// @Summary Envia un correo electronico
// @Description Envia un correo electronico, recive un email para poder hacer esto
// @Tags Email
// @Accept json
// @Produce json
// @Param  emailTo body types.EmailReciber true "a quien se le enviara el correo"
// @Success 200 {object} map[string]string "Mensaje de que salio bien"
// @Failure 400 {object} response.Response "Error en el body"
// @Failure 500 {object} response.Response "Error del servidor"
// @Router /sendEmail [post]
func SendEmail(storage storage.Storage, email utils.EmailHelperSmtp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var emailto types.EmailReciber
		err := json.NewDecoder(r.Body).Decode(&emailto)
		if err != nil {
			slog.Error("error en json decoder")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body %v", err)))
			return
		}

		if err = validator.New().Struct(emailto); err != nil {
			slog.Error("error de validación")
			validateErros := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErros))
			return
		}

		emailT, err := storage.GetEmailTemplate("notificacion_generica")
		slog.Info("Obteniendo template generica de notifiacion")
		if err != nil {
			slog.Error("Error al intentar obtener la template")
			if err.Error() == "no student found" {
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
				return
			}

			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		emailData := types.EmailData{
			Subject:       "Registration Confirmation",
			Greeting:      "Hello John Doe,",
			Message:       "Thank you for registering for the upcoming course! We're excited to have you onboard.",
			ActionDetails: "Your registration was successfully completed. Please check your inbox for further instructions.",
			NextSteps: []string{
				"1. Check your email for confirmation and schedule details.",
				"2. Prepare any required documents for the first session.",
				"3. Get ready to start learning!",
			},
			RecipientName: "John Doe",
			SenderName:    "The School Administration",
			SenderTitle:   "Your School Team",
		}

		tmplt, err := template.New("email").Parse(emailT)
		if err != nil {
			slog.Error("Error al intentar renderizar la template")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		var body bytes.Buffer
		err = tmplt.Execute(&body, emailData)
		if err != nil {
			slog.Error("Error al intentar poner los datos en el template")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		err = email.SimpleEmailSend("Noticias importantes", body.String(), emailto.To)
		if err != nil {
			slog.Error("Error al enviar email")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Error("Se ha enviado correctamente el correo electronico")
		response.WriteJson(w, http.StatusOK, map[string]string{"succes": "OK"})

	}
}
