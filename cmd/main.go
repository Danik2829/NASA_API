package main

import (
	// "github.com/gofiber/fiber/v3"
	"fmt"
	"log"

	"github.com/NASA_API/internal/config"
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

	fmt.Printf("Success: %#v\n %#v\n", metadataStore, imageStore)
}
