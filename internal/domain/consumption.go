package domain

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
)

var (
	errInvalidTimeFormat = errors.New("invalid time format")
)

type DayConsumption struct {
	Day              string
	Date             string
	Consumption      []float64
	TotalConsumption float64
}

func (d DayConsumption) GetTime() (time.Time, error) {
	s := strings.Split(d.Date, "/")
	if len(s) != 3 {
		return time.Time{}, errInvalidTimeFormat
	}

	year, err := strconv.Atoi(s[0])
	if err != nil {
		return time.Time{}, errInvalidTimeFormat
	}

	month, err := strconv.Atoi(s[1])
	if err != nil {
		return time.Time{}, errInvalidTimeFormat
	}

	day, err := strconv.Atoi(s[2])
	if err != nil {
		return time.Time{}, errInvalidTimeFormat
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil

}

type MonthConsumption struct {
	Month       string
	Consumption []DayConsumption
}

type PriceDuration struct {
	SEKPerkWh    float64
	EURPerkWh    float64
	ExchangeRate float64
	TimeStart    time.Time
	TimeEnd      time.Time
}

type HourConsumptionWithPrice struct {
	Start       time.Time
	End         time.Time
	SEKPerkWH   float64
	EURPerkWH   float64
	Consumption float64
}

func (h HourConsumptionWithPrice) GetTotalPriceSEK() float64 {
	return h.Consumption * h.SEKPerkWH
}

func (h HourConsumptionWithPrice) GetTotalPriceEUR() float64 {
	return h.Consumption * h.EURPerkWH
}

type DayConsumptionWithPrice struct {
	Date            time.Time
	HourConsumption []HourConsumptionWithPrice
}

func (d DayConsumptionWithPrice) GetTotalPriceSEK() float64 {
	var total float64
	for _, h := range d.HourConsumption {
		total += h.GetTotalPriceSEK()
	}
	return total
}

func (d DayConsumptionWithPrice) GetTotalPriceEUR() float64 {
	var total float64
	for _, h := range d.HourConsumption {
		total += h.GetTotalPriceEUR()
	}
	return total
}

func (d DayConsumptionWithPrice) GetTotalConsumption() float64 {
	var total float64
	for _, h := range d.HourConsumption {
		total += h.Consumption
	}
	return total
}

type Region int64

const (
	SE1 Region = iota
	SE2
	SE3
	SE4
)

func (r Region) String() string {
	switch r {
	case SE1:
		return "SE1"
	case SE2:
		return "SE2"
	case SE3:
		return "SE3"
	case SE4:
		return "SE4"
	default:
		return "unknown"
	}
}
