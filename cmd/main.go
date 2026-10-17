package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/NASA_API/internal/config"
	"github.com/NASA_API/internal/handler"
	"github.com/NASA_API/internal/nasa"
	"github.com/NASA_API/internal/service"
	"github.com/NASA_API/internal/store"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	metadataStore, err := store.NewMetaDataStore(conf.DB)
	if err != nil {
		log.Fatal(err)
	}

	imageStore, err := store.NewMinioStore(conf.MinioConfig)
	if err != nil {
		log.Fatal(err)
	}

	nasaClient := nasa.NewNasaApi()
	service := service.NewService(metadataStore, imageStore, nasaClient)
	handler := handler.NewHandler(service)

	server := fiber.New()
	handler.Mount(server)

	if err := server.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
