package main

import (
	"context"
	"log"

	"app/internal/config"
	"app/internal/db"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	pool, err := db.NewPool(context.Background(), cfg.DB.DSN())
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	convRepo := repository.NewConversationRepository(pool)
	msgRepo := repository.NewMessageRepository(pool)

	// Services
	userSvc := service.NewUserService(userRepo)
	convSvc := service.NewConversationService(convRepo)
	msgSvc := service.NewMessageService(msgRepo)

	// Handlers
	userH := handler.NewUserHandler(userSvc)
	convH := handler.NewConversationHandler(convSvc)
	msgH := handler.NewMessageHandler(msgSvc)
	pubH := handler.NewPublicHandler(msgSvc)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	// Users
	app.Post("/users", userH.Create)
	app.Get("/users", userH.List)
	app.Get("/users/:id", userH.GetByID)

	// Conversations
	app.Post("/conversations", convH.GetOrCreate)
	app.Get("/conversations/:id", convH.GetByID)
	app.Get("/users/:userId/conversations", convH.ListByUser)

	// Messages
	app.Post("/messages", msgH.Send)
	app.Get("/messages/:conversationId", msgH.ListByConversation)

	// Public discover
	app.Get("/public/messages", pubH.ListMessages)

	log.Printf("listening on :%s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
