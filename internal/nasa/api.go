package nasa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NASA_API/internal/core"
)

const NasaUrl = "https://science.nasa.gov/wp-json/wp/v2/apod-basic/%s"

type NasaApi struct {
}

func NewNasaApi() *NasaApi {
	return &NasaApi{}
}

func (n *NasaApi) GetApodByDate(date string) (*core.Apod, error) {
	s, _ := time.Parse("2006-01-02", date)
	out := s.Format("060102")
	url := fmt.Sprintf(NasaUrl, out)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error while request api.nasa.gov: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"NASA API returned %s: %s",
			resp.Status,
			string(body),
		)
	}

	apod := core.Apod{}
	err = json.NewDecoder(resp.Body).Decode(&apod)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing json: %w", err)
	}

	return &apod, nil
}
