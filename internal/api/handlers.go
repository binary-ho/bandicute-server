package api

import (
	"bandicute-server/internal/service"
	"bandicute-server/pkg/logger"
	"encoding/json"
	"net/http"
	"strings"
)

type Application struct {
	writer *service.Writer
}

func NewApplication(writer *service.Writer) *Application {
	return &Application{
		writer: writer,
	}
}

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (app *Application) WriteAllMembersPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	app.writer.WriteAllMembersPost(r.Context())
	app.sendResponse(w, "Started writing all members' posts")
}

func (app *Application) WriteByStudy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studyId := strings.TrimPrefix(r.URL.Path, "/studies/")
	if studyId == "" {
		app.sendError(w, "studyId is required", http.StatusBadRequest)
		return
	}

	app.writer.WriteByStudy(r.Context(), studyId)
	app.sendResponse(w, "Started writing posts for study")
}

func (app *Application) WriteByMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	memberId := strings.TrimPrefix(r.URL.Path, "/members/")
	if memberId == "" {
		app.sendError(w, "memberId is required", http.StatusBadRequest)
		return
	}

	app.writer.WriteByMember(r.Context(), memberId)
	app.sendResponse(w, "Started writing posts for member")
}

func (app *Application) sendResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", logger.Fields{
			"error": err.Error(),
		})
	}
}

func (app *Application) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{Error: message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode error response", logger.Fields{
			"error": err.Error(),
		})
	}
}
