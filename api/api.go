package api

import (
	"database/sql"

	"github.com/froggy-12/mooshroombase/config"
	"github.com/froggy-12/mooshroombase/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	addr          string
	mongoClient   *mongo.Client
	redisClient   *redis.Client
	mariadbClient *sql.DB
}

func APIServer(addr string, mongoClient *mongo.Client, redisClient *redis.Client, mariadbClient *sql.DB) *Server {
	return &Server{
		addr:          addr,
		mongoClient:   mongoClient,
		redisClient:   redisClient,
		mariadbClient: mariadbClient,
	}
}

func (s *Server) Start() error {
	app := fiber.New(fiber.Config{
		BodyLimit: config.Configs.BodySizeLimit,
	})

	// subrouter groups
	featuredRouter := app.Group("/api/v1")
	fileUpload := app.Group("/api/file-upload/v1")

	// Featured Routes
	routes.FeaturedRoutes(featuredRouter)
	routes.FileStorageRoutes(fileUpload)

	return app.Listen(s.addr)
}
