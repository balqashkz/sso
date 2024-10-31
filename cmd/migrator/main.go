package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	// Библиотека для миграций
	"github.com/golang-migrate/migrate/v4"

	// Драйвер для миграций SQLite 3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to store the migrations")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to store the migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Path to store the migrations")

	flag.Parse()

	if storagePath == "" {
		panic("Missing required flag: -storage-path")
	}
	if migrationsPath == "" {
		panic("Missing required flag: -migrations-path")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("No new migrations to apply.")
	} else {
		fmt.Println("Migrations successfully applied.")
	}
}
