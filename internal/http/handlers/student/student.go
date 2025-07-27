package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/prakash8999/go_rest_apis/internal/types"
	"github.com/prakash8999/go_rest_apis/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Empty request body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusAccepted, response.GeneralError(err))
			return
		}

		//validate request
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidatorError(validateErrs))
			return
		}

		response.WriteJson(w, http.StatusAccepted, map[string]string{"success": "ok"})
		// w.Write([]byte("Welcome to the student API"))

	}
}
