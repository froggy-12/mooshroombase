package routes

import (
	"github.com/froggy-12/mooshroombase/internal/auth/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func MongoAuthRoutes(router fiber.Router, mongoClient *mongo.Client) {
	validate := validator.New()
	router.Post("/create_user_email_password", func(c *fiber.Ctx) error {
		err := handlers.CreateUserWithEmailAndPassword(c, mongoClient, *validate)
		return err
	})
	router.Post("/log_in_email_password", func(c *fiber.Ctx) error {
		err := handlers.LogInWithEmailAndPassword(c, mongoClient, *validate)
		return err
	})
	router.Get("/email_verified", func(c *fiber.Ctx) error {
		err := handlers.VerifyEmail(c, mongoClient, *validate)
		return err
	})
}
