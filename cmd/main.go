package main

import (
	"fmt"
	"github.com/instinctG/tender/internal/db"
	"github.com/instinctG/tender/internal/server"
	"github.com/instinctG/tender/internal/service"
	"github.com/joho/godotenv"
	"log"
)

// Run инициализирует и запускает приложение, устанавливает соединение с базой данных,
// выполняет миграции базы данных, создает сервис статистики и запускает HTTP сервер.
func Run() error {
	fmt.Println("starting up our application")

	database, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return err
	}

	if err = database.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	tenderService := service.NewService(database)

	httpHandler := server.NewHandler(tenderService)
	if err = httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

// main является точкой входа в приложение, вызывает функцию Run и обрабатывает возможные ошибки.
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("err loading .env file : %v", err)
	}
	fmt.Println("Running tender service")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
