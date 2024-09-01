package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/joshuandeleva/go-ticket-backend/config"
	"github.com/joshuandeleva/go-ticket-backend/db"
	"github.com/joshuandeleva/go-ticket-backend/handlers"
	"github.com/joshuandeleva/go-ticket-backend/middlewares"
	"github.com/joshuandeleva/go-ticket-backend/repositories"
	"github.com/joshuandeleva/go-ticket-backend/services"
)

func main(){
	enConfig := config.NewEnvConfig()
	db := db.Init(enConfig , db.DBMigrator)
	app := fiber.New(fiber.Config{
		AppName: "Go Ticket API",
		Prefork: true,
		ServerHeader: "Go Ticket API",
	})
	// repository
	eventRepository := repositories.NewEventRespository(db)
	ticketRepository := repositories.NewTicketRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	// service

	authService := services.NewAuthService(authRepository)

	// routes
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	privateRoutes := server.Use(middlewares.AuthProtected(db))

	//handlers

	handlers.NewEventHandler(privateRoutes.Group("/events"), eventRepository)
	handlers.NewTicketHandler(privateRoutes.Group("/tickets"), ticketRepository)

	app.Listen(fmt.Sprintf(":%s", enConfig.ServerPort))

}