package main

import (
	"SkillsRock/internal/api/handlers"
	"SkillsRock/internal/api/router"
	"SkillsRock/internal/config"
	"SkillsRock/internal/databases/postgres"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

func runMigrations(dsn string) error {
	// Убедитесь, что путь к миграциям верный
	migrationsPath := "file://internal/databases/postgres/migrations"

	m, err := migrate.New(
		migrationsPath,
		"pgx"+dsn[8:],
	)
	if err != nil {
		return fmt.Errorf("[ERROR] не удалось создать экземпляр миграции: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("[ERROR] не удалось запустить миграции: %w", err)
	}

	log.Println("Миграции успешно применены")
	return nil
}
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("ОшИбка загрузки конфига: %v", err)
	}
	// Подключение к базе данных с помощью pgx
	conn, err := pgx.Connect(context.Background(), cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}
	defer conn.Close(context.Background())

	// Проверка соединения
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Не удалось выполнить ping базы данных: %v", err)
	}

	// Запуск миграций
	if err := runMigrations(cfg.Database.GetDSN()); err != nil {
		log.Fatalf(" %v", err)
	}

	log.Println("Миграции успешно выполнены")
	taskdb := postgres.NewTaskDb(conn)

	taskHandler := handlers.NewTaskHandler(taskdb)
	router := router.SetupRouter(*taskHandler)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		router.Listen(":8000")
	}()
	wg.Wait()
}
