package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
