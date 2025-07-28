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
	"github.com/prakash8999/go_rest_apis/internal/storage"
	"github.com/prakash8999/go_rest_apis/internal/types"
	"github.com/prakash8999/go_rest_apis/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("User created with id", slog.String("userId ", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
		// w.Write([]byte("Welcome to the student API"))

	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		//parse  id to int64
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("error parsing id", slog.String("id", id))

			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			// slog.Error()
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, 200, student)
	}
}
