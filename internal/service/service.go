package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/NASA_API/internal/core"
)

type Service struct {
	metaDataStore core.MetaDataActions
	imageStore    core.ImageActions
	nasaAPI       core.NasaAPI
}

func NewService(m core.MetaDataActions, i core.ImageActions, api core.NasaAPI) *Service {
	return &Service{
		metaDataStore: m,
		imageStore:    i,
		nasaAPI:       api,
	}
}

func (s *Service) GetAllApods() ([]core.Apod, error) {
	log.Println("service: GetAllApods")
	apods, err := s.metaDataStore.GetAllApods()
	if err != nil {
		log.Println("service: GetAllApods error:", err)
		return nil, err
	}
	return apods, nil
}

func (s *Service) GetByDate(date string, ctx context.Context) (*core.Apod, error) {
	log.Println("service: GetByDate", date)

	apod, err := s.metaDataStore.GetApodByDate(date)
	if err != nil {
		log.Println("service: GetByDate metadata lookup error:", err)
		return nil, err
	}
	if apod != nil {
		log.Println("service: GetByDate cache hit", date)
		return apod, nil
	}

	log.Println("service: GetByDate cache miss, calling NASA API", date)
	apod, err = s.nasaAPI.GetApodByDate(date)
	if err != nil {
		log.Println("service: GetByDate NASA API error:", err)
		return nil, err
	}

	if apod.MediaType == "image" {
		log.Println("service: GetByDate downloading image", apod.HdURL)
		resp, err := http.Get(apod.HdURL)
		if err != nil {
			log.Println("service: GetByDate image download error:", err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("service: GetByDate image download failed:", resp.Status)
			return nil, fmt.Errorf("NASA image download failed: %s", resp.Status)
		}

		objectName := "apod/" + apod.Date + ".jpg"
		log.Println("service: GetByDate storing image", objectName)

		err = s.imageStore.PutImage(
			ctx,
			objectName,
			resp.Body,
			resp.ContentLength,
			resp.Header.Get("Content-Type"),
		)
		if err != nil {
			log.Println("service: GetByDate put image error:", err)
			return nil, err
		}
	}

	log.Println("service: GetByDate saving apod", apod.Date)
	if err := s.metaDataStore.SaveApod(apod); err != nil {
		log.Println("service: GetByDate save apod error:", err)
		return nil, err
	}

	return apod, nil
}

func (s *Service) GetByDateImage(date string, ctx context.Context) ([]byte, error) {
	log.Println("service: GetByDateImage", date)

	apod, err := s.metaDataStore.GetApodByDate(date)
	if err != nil {
		log.Println("service: GetByDateImage metadata lookup error:", err)
		return nil, err
	}
	if apod != nil {
		if apod.MediaType != "image" {
			log.Println("service: GetByDateImage not an image", date)
			return nil, fmt.Errorf("APOD for %s is not an image", date)
		}

		objectName := "apod/" + apod.Date + ".jpg"
		log.Println("service: GetByDateImage reading from store", objectName)
		image, err := s.imageStore.GetImage(ctx, objectName)
		if err != nil {
			log.Println("service: GetByDateImage get image error:", err)
			return nil, err
		}

		return image, nil
	}

	log.Println("service: GetByDateImage cache miss, calling NASA API", date)
	apod, err = s.nasaAPI.GetApodByDate(date)
	if err != nil {
		log.Println("service: GetByDateImage NASA API error:", err)
		return nil, err
	}

	log.Println("service: GetByDateImage downloading image", apod.HdURL)
	resp, err := http.Get(apod.HdURL)
	if err != nil {
		log.Println("service: GetByDateImage image download error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("service: GetByDateImage image download failed:", resp.Status)
		return nil, fmt.Errorf("NASA image download failed: %s", resp.Status)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("service: GetByDateImage read body error:", err)
		return nil, err
	}
	log.Println(apod.Date)
	objectName := "apod/" + apod.Date + ".jpg"

	log.Println("service: GetByDateImage storing image", objectName)
	err = s.imageStore.PutImage(
		ctx,
		objectName,
		bytes.NewReader(imageBytes),
		int64(len(imageBytes)),
		resp.Header.Get("Content-Type"),
	)
	if err != nil {
		log.Println("service: GetByDateImage put image error:", err)
		return nil, err
	}

	log.Println("service: GetByDateImage saving apod", apod.Date)
	if err := s.metaDataStore.SaveApod(apod); err != nil {
		log.Println("service: GetByDateImage save apod error:", err)
		return nil, err
	}

	return imageBytes, nil
}
