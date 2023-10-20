package main

import (
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(500).SendString("Upload failed")
		}

		files := form.File["files"]

		for _, file := range files {
			// Open the file
			uploadedFile, err := file.Open()
			if err != nil {
				return c.Status(500).SendString("Error opening file")
			}
			defer uploadedFile.Close()

			// Create a new file on the server to write chunks
			newFile, err := os.Create("converted_.mov")
			if err != nil {
				return c.Status(500).SendString("Error creating file")
			}
			defer newFile.Close()

			// Copy chunks from uploaded file to the new file
			if _, err := io.Copy(newFile, uploadedFile); err != nil {
				return c.Status(500).SendString("Error copying file")
			}
		}

		return c.SendString("File(s) uploaded successfully")
	})

	app.Listen(":3000")
}
