package api

import "github.com/gofiber/fiber/v2"

func (app *Application) Routes() *fiber.App {
	router := fiber.New()

	router.Post("/posts", app.WriteAllMembersPost)
	router.Post("/studies/:studyId/posts", app.WriteByStudy)
	router.Post("/members/:memberId/posts", app.WriteByMember)
	return router
}
