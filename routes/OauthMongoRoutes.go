package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/froggy-12/mooshroombase/config"
	"github.com/froggy-12/mooshroombase/internal/types"
	"github.com/froggy-12/mooshroombase/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"go.mongodb.org/mongo-driver/mongo"
)

func OAuthMongoRoutes(router fiber.Router, mongoClient *mongo.Client) {
	goth.UseProviders(
		github.New(config.Configs.GithubKey, config.Configs.GithubSecret, config.Configs.Back_End_URL+"/api/mongo/oauth/v1/callback/github"),
		google.New(config.Configs.GoogleKey, config.Configs.GoogleSecret, config.Configs.Back_End_URL+"/api/mongo/oauth/v1/callback/google"),
	)

	router.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	router.Get("/callback/:provider", func(c *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(c)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Something Went Wrong: " + err.Error()})
		}

		newUser := types.User_Mongo_Oauth{
			ID:             user.UserID,
			UserName:       user.Name,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			ProfilePicture: user.AvatarURL,
			Email:          user.Email,
			Verified:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if c.Params("provider") == "github" {
			coll := mongoClient.Database("mooshroombase").Collection("github_oauth_users")
			_, err := coll.InsertOne(context.Background(), newUser)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(types.ErrorResponse{Error: "Failed to create user in database: " + err.Error()})
			}
		}
		if c.Params("provider") == "google" {
			coll := mongoClient.Database("mooshroombase").Collection("google_oauth_users")
			_, err = coll.InsertOne(context.Background(), newUser)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(types.ErrorResponse{Error: "Failed to create user in database: " + err.Error()})
			}
		}

		token, err := utils.GenerateOauthJWTToken(newUser.ID, config.Configs.JWTTokenExpiration, config.Configs.JWTSecret)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Failed to generate jwt token"})
		}

		cookie := &fiber.Cookie{
			Name:     "jwtToken",
			Value:    token,
			Path:     "/",
			HTTPOnly: true,
			MaxAge:   config.Configs.JWTCookieAge,
			Secure:   true,
		}

		c.Cookie(cookie)
		return c.Status(http.StatusOK).JSON(types.HttpSuccessResponse{Message: "User Has been created"})
	})

}
