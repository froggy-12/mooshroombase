package handlers

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/froggy-12/mooshroombase/config"
	"github.com/froggy-12/mooshroombase/internal/types"
	"github.com/froggy-12/mooshroombase/smtp_configs"
	"github.com/froggy-12/mooshroombase/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserWithEmailAndPassword(c *fiber.Ctx, mongoClient *mongo.Client, validate validator.Validate) error {
	jwtTokenCookie := c.Cookies("jwtToken")
	if jwtTokenCookie != "" {
		return c.Status(http.StatusBadGateway).JSON(types.ErrorResponse{Error: "we detected you have jwtToken you cant sign up either your jwtToken is invalid or its outdated try to log in"})
	}

	collection := mongoClient.Database("mooshroombase").Collection("users")
	var user types.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Invalid Request body"})
	}

	if err := validate.Struct(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Failed to Hash the password: " + err.Error()})
	}

	Id := uuid.New().String()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	verificationToken := rand.Intn(90000) + 100000

	newUser := types.User_Mongo{
		FirstName:         user.FirstName,
		ID:                Id,
		LastName:          user.LastName,
		Email:             user.Email,
		UserName:          user.UserName,
		Password:          string(hashedPassword),
		ProfilePicture:    config.Configs.DefaultProfilePicURL,
		Verified:          false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		VerificationToken: verificationToken,
	}

	_, err = collection.InsertOne(context.Background(), newUser)

	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(types.ErrorResponse{Error: "failed to create new user into the database: " + err.Error()})
	}

	if config.Configs.EmailVerificationAllowed {
		err = smtp_configs.SendVerificationEmail(newUser.Email, newUser.VerificationToken)

		if err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(types.ErrorResponse{Error: "Failed to sent email and aborted creating user: " + err.Error()})
		}
	}

	return c.Status(http.StatusCreated).JSON(types.HttpSuccessResponse{Message: "User Has been created to the database make sure to verify user's email"})
}

func LogInWithEmailAndPassword(c *fiber.Ctx, mongoClient *mongo.Client, validate validator.Validate) error {
	tokenString := c.Cookies("jwtToken")
	coll := mongoClient.Database("mooshroombase").Collection("users")
	if tokenString != "" {
		userId, expired, err := utils.ReadJWTToken(tokenString, config.Configs.JWTSecret)
		if err != nil || expired {
			return utils.LogIn(c, coll, validate, config.Configs.JWTTokenExpiration, config.Configs.JWTSecret, config.Configs.JWTCookieAge)
		}
		_, err = utils.FindUserFromMongoDBUsingID(userId, coll)
		if err == nil {
			return c.Status(http.StatusAlreadyReported).JSON(types.HttpSuccessResponse{Message: "You are already logged in"})
		}
	}

	return utils.LogIn(c, coll, validate, config.Configs.JWTTokenExpiration, config.Configs.JWTSecret, config.Configs.JWTCookieAge)

}

func VerifyEmail(c *fiber.Ctx, mongoClient *mongo.Client, validate validator.Validate) error {
	coll := mongoClient.Database("mooshroombase").Collection("users")

	email := c.Query("email")
	code := c.Query("code")

	if email == "" || code == "" {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Email and code are required"})
	}

	if err := validate.Var(email, "required,email"); err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Invalid email: " + err.Error()})
	}

	codeInt, err := strconv.Atoi(code)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Invalid code: " + err.Error()})
	}

	user, err := utils.FindUserFromMongoDBUsingEmail(email, coll)
	if err != nil {
		return c.Status(http.StatusBadGateway).JSON(types.ErrorResponse{Error: "User not Found"})
	}

	if user.VerificationToken == codeInt {
		_, err := coll.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"verified": true}})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(types.ErrorResponse{Error: "Failed to update user verification status"})
		}
		return c.Status(http.StatusOK).JSON(types.HttpSuccessResponse{Message: "Email verified successfully"})
	} else {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Wrong Code Provided"})
	}
}
