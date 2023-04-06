package elpriserjustnu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GLindebyV/electricityconsumption/internal/domain"
	"github.com/go-errors/errors"
)

const (
	elpriserjustnuEndpoint = "https://www.elprisetjustnu.se/api"
	version1               = "v1"
)

type Client struct {
	baseURL    string
	version    string
	httpClient *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		baseURL:    elpriserjustnuEndpoint,
		version:    version1,
		httpClient: client,
	}
}

func (c *Client) GetPricePoints(t time.Time, region domain.Region) ([]domain.PriceDuration, error) {
	month := fmt.Sprintf("%d", t.Month())
	if len(month) == 1 {
		month = fmt.Sprintf("0%s", month)
	}
	day := fmt.Sprintf("%d", t.Day())
	if len(day) == 1 {
		day = fmt.Sprintf("0%s", day)
	}

	date := fmt.Sprintf("%d/%s-%s", t.Year(), month, day)
	url := fmt.Sprintf("%s/%s/prices/%s_%s.json", c.baseURL, c.version, date, region.String())
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var points pricePoints
	if err := json.NewDecoder(resp.Body).Decode(&points); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	priceDurations, err := points.ToPriceDuration()
	if err != nil {
		return nil, err
	}

	return priceDurations, nil
}
