package api

import (
	"bandicute-server/config"
	"bandicute-server/internal/job"
	"bandicute-server/internal/service"
	"bandicute-server/pkg/logger"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type Application struct {
	config            *config.Config
	writer            *service.Writer
	serviceDispatcher *service.Dispatcher
	scheduler         *job.Scheduler
}

func NewApplication(config *config.Config,
	writer *service.Writer,
	dispatcher *service.Dispatcher,
	scheduler *job.Scheduler) *Application {
	return &Application{
		config:            config,
		writer:            writer,
		serviceDispatcher: dispatcher,
		scheduler:         scheduler,
	}
}

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (app *Application) Run() {
	app.scheduler.Start()
	defer app.scheduler.Shutdown()

	app.serviceDispatcher.Run()

	fiberApp := app.Routes()

	err := fiberApp.Listen(getStringPort(app.config.Server.Port))
	if err != nil {
		logger.Fatal("Server Error", logger.Fields{
			"error": err.Error(),
		})
	}
}

func (app *Application) WriteAllMembersPost(c *fiber.Ctx) error {
	app.writer.WriteAllMembersPost(c.Context())
	return c.JSON(fiber.Map{
		"message": "Started writing all members' posts",
	})
}

func (app *Application) WriteByStudy(c *fiber.Ctx) error {
	studyId := c.Params("studyId")
	if studyId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "studyId is required",
		})
	}

	app.writer.WriteByStudy(c.Context(), studyId)
	return c.JSON(fiber.Map{
		"message": "Started writing posts for study",
	})
}

func (app *Application) WriteByMember(c *fiber.Ctx) error {
	memberId := c.Params("memberId")
	if memberId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "memberId is required",
		})
	}

	app.writer.WriteByMember(c.Context(), memberId)
	return c.JSON(fiber.Map{
		"message": "Started writing posts for member",
	})
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

func getStringPort(port int) string {
	return ":" + strconv.Itoa(port)
}
