package api

import (
	"database/sql"
	"strings"

	"github.com/froggy-12/mooshroombase/config"
	"github.com/froggy-12/mooshroombase/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	corsMiddleWare := cors.New(cors.Config{
		AllowOrigins: strings.Join(config.Configs.AllowedCorsOrigin, ", "),
		AllowMethods: strings.Join([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}, ", "),
		AllowHeaders: "*",
		MaxAge:       config.Configs.CorsHeadersMaxAge,
	})

	// subrouter groups
	featuredRouter := app.Group("/api/v1")
	fileUpload := app.Group("/api/file-upload/v1", corsMiddleWare)

	// Featured Routes
	routes.FeaturedRoutes(featuredRouter)
	routes.FileStorageRoutes(fileUpload)

	// mongo auth routes
	if config.Configs.Authentication {
		if config.Configs.PrimaryDB == "mongodb" {
			mongoAuthRouter := app.Group("/api/mongo/auth/v1", corsMiddleWare)
			mongoOAuthRouter := app.Group("/api/mongo/oauth/v1", corsMiddleWare)
			routes.MongoAuthRoutes(mongoAuthRouter, s.mongoClient)
			routes.OAuthMongoRoutes(mongoOAuthRouter, s.mongoClient)
		}
	}

	return app.Listen(s.addr)
}
