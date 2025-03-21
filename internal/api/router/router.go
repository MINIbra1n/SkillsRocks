package router

import (
	"SkillsRock/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(handlers handlers.TaskHandler) *fiber.App {
	r := fiber.New()

	task := r.Group("/task")
	task.Post("", handlers.AddTasks)
	task.Get("", handlers.GetTask)
	task.Put("/:id", handlers.UpdateTask)
	task.Delete("/:id", handlers.DeleteTaskWithId)

	return r

}
