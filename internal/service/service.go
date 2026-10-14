package service

import (
	"context"
	"fmt"
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

func (s *Service) GetByDate(date string, ctx context.Context) (*core.Apod, error) {
	apod, err := s.metaDataStore.GetApodByDate(date)
	if err != nil {
		return nil, err
	}
	if apod != nil {
		fmt.Println("It was from cache))))")
		return apod, nil
	}

	apod, err = s.nasaAPI.GetApodByDate(date)
	if err != nil {
		return nil, err
	}

	if apod.MediaType == "image" {
		resp, err := http.Get(apod.HdURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("NASA image download failed: %s", resp.Status)
		}

		objectName := "apod/" + apod.Date + ".jpg"

		err = s.imageStore.PutImage(
			ctx,
			objectName,
			resp.Body,
			resp.ContentLength,
			resp.Header.Get("Content-Type"),
		)
		if err != nil {
			return nil, err
		}
	}

	if err := s.metaDataStore.SaveApod(apod); err != nil {
		return nil, err
	}

	return apod, nil
}
