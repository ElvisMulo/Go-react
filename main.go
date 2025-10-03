package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        string `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello")

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	// Get all todos
	app.Get("/api/todos/all", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Create a to do
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := Todo{}

		if err := c.BodyParser(&todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "Body is required"})
		}

		todo.ID = strconv.Itoa(len(todos) + 1)
		todos = append(todos, todo)

		return c.Status(201).JSON(todo)
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
	})

	//Delete a to do
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"msg": "Todo deleted"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT)) // if there is any error
}
