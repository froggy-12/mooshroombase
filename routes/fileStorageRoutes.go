package routes

import (
	"github.com/froggy-12/mooshroombase/storage"
	"github.com/gofiber/fiber/v2"
)

func FileStorageRoutes(router fiber.Router) {
	router.Post("/upload/image/single", storage.HandleUploadImageFile)
	router.Post("/upload/image/multi", storage.HandleUploadMultipleImageFile)
	router.Post("/upload/music/single", storage.HandleUploadSingleMusicFile)
	router.Post("/upload/music/multi", storage.HandleUploadMultipleMusicFile)
	router.Post("/upload/video/single", storage.HandleUploadSingleVideoFile)
	router.Post("/upload/video/multi", storage.HandleUploadMultiVideoFile)
	router.Post("/upload/video/multi", storage.HandleUploadMultiVideoFile)
	router.Post("/upload/any/single", storage.HandleAnyFormatSingleFile)
	router.Post("/upload/any/multi", storage.HandleAnyFormatMultiFile)
	router.Delete("/deletefile", storage.HandleDeleteFile)
}
