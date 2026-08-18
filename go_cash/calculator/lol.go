package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var todos = []Todo{
	{ID: 1, Task: "Learn Go", Done: false},
	{ID: 2, Task: "Build a Fiber server", Done: true},
}

func main() {
	app := fiber.New()

	// GET /
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Fiber server is running!")
	})

	// GET /todos - list all todos
	app.Get("/todos", func(c *fiber.Ctx) error {
		return c.JSON(todos)
	})

	// GET /todos/:id - get a single todo by ID
	app.Get("/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid id",
			})
		}

		for _, t := range todos {
			if t.ID == id {
				return c.JSON(t)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "todo not found",
		})
	})

	// POST /todos - create a new todo
	app.Post("/todos", func(c *fiber.Ctx) error {
		var newTodo Todo
		if err := c.BodyParser(&newTodo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		newTodo.ID = len(todos) + 1
		todos = append(todos, newTodo)

		return c.Status(fiber.StatusCreated).JSON(newTodo)
	})

	log.Fatal(app.Listen(":3000"))
}