package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/froggy-12/mooshroombase/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// utility function for debug logging
func DebugLogger(info, message any) {
	fmt.Printf("[debug] [%v]: %v \n", info, message)
}

func IsImage(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".avif", ".jpig"}
	for _, ext := range extensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}
	return false
}

func IsMusic(filename string) bool {
	extensions := []string{".mp3", ".wav", ".ogg", ".flac", ".m4a"}
	for _, ext := range extensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}
	return false
}

func IsVideo(filename string) bool {
	ext := filepath.Ext(filename)
	ext = strings.ToLower(ext)
	return ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".wmv"
}

func UploadFile(_ *fiber.Ctx, file *multipart.FileHeader, folder string) (string, error) {
	// Generate a unique filename
	uuid := uuid.New()
	filename := fmt.Sprintf("%s%s", uuid, filepath.Ext(file.Filename))

	// Save the file to the uploads folder
	uploadsDir := filepath.Join(".", "uploads", folder)
	err := os.MkdirAll(uploadsDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	filepath := filepath.Join(uploadsDir, filename)
	fileStream, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileStream.Close()

	// Create the file
	f, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Copy the file contents
	_, err = io.Copy(f, fileStream)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func DeleteFile(filename, folder string) error {
	filepath := filepath.Join(".", "uploads", folder, filename)
	return os.Remove(filepath)
}

func UploadAnyFile(_ *fiber.Ctx, file *multipart.FileHeader, folder, filename string) error {
	// Save the file to the uploads folder
	uploadsDir := filepath.Join(".", "uploads", folder)
	err := os.MkdirAll(uploadsDir, os.ModePerm)
	if err != nil {
		return err
	}

	filepath := filepath.Join(uploadsDir, filename)
	fileStream, err := file.Open()
	if err != nil {
		return err
	}
	defer fileStream.Close()

	// Create the file
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the file contents
	_, err = io.Copy(f, fileStream)
	if err != nil {
		return err
	}

	return nil
}

func GenerateJWTToken(id string, jwtExpirationTime time.Time, jwtSecret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id,
		"expr": jwtExpirationTime.Unix(),
		"iat":  time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func SetJWTHttpCookies(c *fiber.Ctx, token string, message string, cookieAge int) error {
	cookie := &fiber.Cookie{
		Name:     "jwtToken",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   cookieAge,
	}
	c.Cookie(cookie)

	return c.Status(http.StatusCreated).JSON(types.HttpSuccessResponse{Message: message})
}

func FindUserFromMongoDBUsingEmail(email string, mongoCollection *mongo.Collection) (types.User_Mongo, error) {
	filter := bson.M{"email": email}
	user := types.User_Mongo{}
	err := mongoCollection.FindOne(context.Background(), filter).Decode(&user)
	return user, err
}

func FindUserFromMongoDBUsingUsername(username string, mongoCollection *mongo.Collection) (types.User_Mongo, error) {
	filter := bson.M{"username": username}
	user := types.User_Mongo{}
	err := mongoCollection.FindOne(context.Background(), filter).Decode(&user)
	return user, err
}

func FindUserFromMongoDBUsingID(id string, mongoCollection *mongo.Collection) (types.User_Mongo, error) {
	filter := bson.M{"id": id}
	user := types.User_Mongo{}
	err := mongoCollection.FindOne(context.Background(), filter).Decode(&user)
	return user, err
}

func ReadJWTToken(token string, jwtSecret string) (string, bool, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", false, err
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return "", false, errors.New("invalid token claims")
	}
	expr, ok := claims["expr"].(float64)
	if !ok {
		return "", false, errors.New("invalid token claims")
	}

	expirationTime := time.Unix(int64(expr), 0)
	if time.Now().After(expirationTime) {
		return "", true, nil
	}

	return userId, false, nil
}

func LogIn(c *fiber.Ctx, coll *mongo.Collection, validate validator.Validate, jwtExpirationTime time.Time, jwtSecret string, cookieAge int) error {
	var details types.LogInDetails
	if err := c.BodyParser(&details); err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Invalid Request body"})
	}

	if err := validate.Struct(&details); err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: err.Error()})
	}

	user, err := FindUserFromMongoDBUsingEmail(details.Email, coll)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "User Doesnt Exist"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(details.Password))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{Error: "Wrong Password Try again"})
	}

	token, err := GenerateJWTToken(user.ID, jwtExpirationTime, jwtSecret)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(types.ErrorResponse{Error: "Failed to generate JWT token"})
	}

	err = SetJWTHttpCookies(c, token, "user with id: "+user.ID+" has been logged in successfully", cookieAge)
	return err

}
