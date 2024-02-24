package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/akatranlp/gs-analysis-go/internal/config"
	"github.com/akatranlp/gs-analysis-go/internal/database"
	"github.com/caarlos0/env/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func main() {
	log.Println("Starting application...")
	log.Println("Loading Application config...")
	_ = godotenv.Load()
	cfg := config.ApplicationConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to get application configuration", err)
	}

	log.Println("Connecting to database...")
	conn, err := sql.Open("sqlite3", "data/db.sqlite3")
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	defer conn.Close()

	log.Println("Migrate database...")
	embedMigrations := database.MigrationFiles()
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal("failed to set dialect", err)
	}

	if err := goose.Up(conn, "migrations"); err != nil {
		log.Fatal("failed to migrate database", err)
	}

	db := database.New(conn)

	log.Println("Starting server...")
	log.Println(os.Getenv("TEST"))

	app := fiber.New()

	app.Get("/api/authors", func(c *fiber.Ctx) error {
		authors, err := db.ListAuthors(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to list authors")
		}
		return c.JSON(authors)
	})

	app.Mount("/", setupFrontend(os.Getenv("FRONTEND_PATH")))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		app.ShutdownWithContext(context.Background())
	}()

	app.Listen(fmt.Sprintf(":%d", 3000))
}

func setupFrontend(frontendPath string) *fiber.App {
	app := fiber.New()
	app.Static("/", frontendPath)

	// Catch all routes
	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := path.Join(frontendPath, "index.html")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })

	return app
}
