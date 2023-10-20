package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/qinains/fastergoding"
)

func main() {
	// Hot reloading
	fastergoding.Run()
	app := fiber.New()

	// App config
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Printf("Logging base route\n")
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

			fmt.Printf("Uploaded file: %v\n", uploadedFile)

			// Create the directory if it doesn't exist
			targetDirectory := "../temp"
			if err := os.MkdirAll(targetDirectory, 0755); err != nil {
				return c.Status(500).SendString(fmt.Sprintf("Error creating file:", err))
			}
			// Create a new file on the server to write chunks
			newFilePath := fmt.Sprintf("../temp/converted_video.mov")
			newFile, err := os.Create(newFilePath)
			if err != nil {
				log.Println("Error creating file:", err)
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
