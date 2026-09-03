package student

import (
	"GO-API/internal/storage"
	"GO-API/internal/types"
	"GO-API/internal/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(stor storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		validate := validator.New()
		if err := validate.Struct(student); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationErrors))
			return
		}

		lastID, err := stor.Createstudent(
			student.Name,
			student.Age,
			student.Email,
		)
		slog.Info("User created successfully", slog.String("userID", fmt.Sprint(lastID)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"ID": fmt.Sprint(lastID)})
	}
}

func GetByID(stor storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for getting a student by ID will go here
		id := r.PathValue("id") // Assuming you have a way to extract the ID from the URL path
		slog.Info("gettint a student", slog.String("userID", id))

		intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("error geting user", slog.String("id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student ID")))
			return
		}

		student, err := stor.GetstudentByID(intid)
		if err != nil {
			slog.Error("error geting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
