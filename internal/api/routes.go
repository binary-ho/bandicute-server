package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (app *Application) Routes() *fiber.App {
	fiberApp := fiber.New()

	app.addCorsSetting(fiberApp)
	app.addRoutes(fiberApp)
	return fiberApp
}

func (app *Application) addCorsSetting(router *fiber.App) fiber.Router {
	corsConfig := app.config.CORS
	return router.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowOrigins,
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
	}))
}

func (app *Application) addRoutes(fiberApp *fiber.App) {
	fiberApp.Post("/posts", app.WriteAllMembersPost)
	fiberApp.Post("/studies/:studyId/posts", app.WriteByStudy)
	fiberApp.Post("/members/:memberId/posts", app.WriteByMember)
}
