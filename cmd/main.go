package main

import (
	// "github.com/gofiber/fiber/v3"
	"context"
	"fmt"
	"log"

	"github.com/NASA_API/internal/config"
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
	_ = service
	fmt.Println(service.GetByDate("261014", context.Background()))
}
