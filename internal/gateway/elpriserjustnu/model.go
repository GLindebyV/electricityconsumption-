package elpriserjustnu

import (
	"time"

	"github.com/GLindebyV/electricityconsumption/internal/domain"
	"github.com/go-errors/errors"
)

type pricePoint struct {
	SEKPerkWh    float64 `json:"SEK_per_kWh"`
	EURPerkWh    float64 `json:"EUR_per_kWh"`
	ExchangeRate float64 `json:"EXR"`
	TimeStart    string  `json:"time_start"`
	TimeEnd      string  `json:"time_end"`
}

type pricePoints []pricePoint

func (p pricePoints) ToPriceDuration() ([]domain.PriceDuration, error) {
	var priceDurations []domain.PriceDuration
	for _, pp := range p {

		timeStart, err := toTime(pp.TimeStart)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}

		timeEnd, err := toTime(pp.TimeEnd)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}

		priceDurations = append(priceDurations, domain.PriceDuration{
			SEKPerkWh:    pp.SEKPerkWh,
			EURPerkWh:    pp.EURPerkWh,
			ExchangeRate: pp.ExchangeRate,
			TimeStart:    timeStart,
			TimeEnd:      timeEnd,
		})
	}
	return priceDurations, nil
}

func toTime(s string) (time.Time, error) {
	layout := "2006-01-02T15:04:05-07:00"
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, errors.Wrap(err, 0)
	}

	return parsedTime, nil
}
