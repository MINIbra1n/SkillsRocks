package handlers

import (
	"SkillsRock/internal/databases/postgres"
	"SkillsRock/internal/domain/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	taskDb postgres.TaskDB
}

func NewTaskHandler(taskDB *postgres.TaskDB) *TaskHandler {
	return &TaskHandler{
		taskDb: *taskDB,
	}
}

func (t *TaskHandler) AddTasks(c *fiber.Ctx) error {
	var req models.Task
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Не удалось разобрать JSON",
		})
	}
	ctx := c.Context()
	id, err := t.taskDb.CreateTask(ctx, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Ошибка при создании задачи",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"id":      id,
		"message": "Задача успешно создана",
	})
}

func (t *TaskHandler) GetTask(c *fiber.Ctx) error {
	var res string
	ctx := c.Context()
	task, _, err := t.taskDb.GetTasks(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Ошибка при создании задачи",
		})
	}
	for _, v := range task {
		res += fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", v.Id, v.Title, v.Description, v.Status, v.Created_at, v.Update_at)
	}

	return c.Status(201).Send([]byte(res))
}

func (t *TaskHandler) DeleteTaskWithId(c *fiber.Ctx) error {
	var id string
	id = c.Params("id")
	d, err := strconv.ParseInt(id, 0, 64)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Ошибка при удалении задачи",
		})
	}
	ctx := c.Context()
	err = t.taskDb.DeleteTaskId(ctx, d)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Ошибка при удалении задачи",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"id":      id,
		"message": "Задача успешно Удалена",
	})
}

func (t *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	var req models.Task
	var id string
	id = c.Params("id")
	d, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Ошибка при изменение задачи",
		})
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Не удалось разобрать JSON",
			"asdsaad": req,
			"id":      id,
		})
	}
	ctx := c.Context()
	err = t.taskDb.UpdateTaskID(ctx, d, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Ошибка при изменении задачи",
			"body":  req,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"id":      id,
		"message": "Задача успешно изменена",
	})
}
