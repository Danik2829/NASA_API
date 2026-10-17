package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	//"strings"

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
	apods, err := s.metaDataStore.GetAllApods()
	if err != nil {
		return nil, err
	}
	return apods, nil
}

func (s *Service) GetByDate(date string, ctx context.Context) (*core.Apod, error) {
	apod, err := s.metaDataStore.GetApodByDate(date)
	if err != nil {
		return nil, err
	}
	if apod != nil {
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

func (s *Service) GetByDateImage(date string, ctx context.Context) ([]byte, error) {
	apod, err := s.metaDataStore.GetApodByDate(date)
	if err != nil {
		return nil, err
	}
	if apod != nil {
		if apod.MediaType != "image" {
			return nil, fmt.Errorf("APOD for %s is not an image", date)
		}

		objectName := "apod/" + apod.Date + ".jpg"
		image, err := s.imageStore.GetImage(ctx, objectName)
		if err != nil {
			return nil, err
		}

		return image, nil
	}

	apod, err = s.nasaAPI.GetApodByDate(date)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(apod.HdURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NASA image download failed: %s", resp.Status)
	}
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(apod.Date)
	objectName := "apod/" + apod.Date + ".jpg"

	err = s.imageStore.PutImage(
		ctx,
		objectName,
		bytes.NewReader(imageBytes),
		int64(len(imageBytes)),
		resp.Header.Get("Content-Type"),
	)
	if err != nil {
		return nil, err
	}

	if err := s.metaDataStore.SaveApod(apod); err != nil {
		return nil, err
	}

	return imageBytes, nil
}
