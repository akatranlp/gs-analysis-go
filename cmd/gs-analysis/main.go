package main

import (
	"database/sql"
	"github.com/akatranlp/gs-analysis-go/internal/config"
	"github.com/akatranlp/gs-analysis-go/internal/database"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"log"
	"os"
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
	db, err := sql.Open("sqlite3", "./db/db.sqlite3")
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	log.Println("Migrate database...")
	embedMigrations := database.MigrationFiles()
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal("failed to set dialect", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("failed to migrate database", err)
	}

	log.Println("Starting server...")
	log.Println(os.Getenv("TEST"))
}
