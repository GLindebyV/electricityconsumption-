package service

import (
	"time"

	"github.com/GLindebyV/electricityconsumption/internal/domain"
)

type ElectricityPrice interface {
	GetPricePoints(t time.Time, region domain.Region) ([]domain.PriceDuration, error)
}

type PricingService struct {
	ElectricityPrice ElectricityPrice
}

func NewPricingService(electricityPrice ElectricityPrice) PricingService {
	return PricingService{
		ElectricityPrice: electricityPrice,
	}
}

func (ps PricingService) GetDayConsumptionWithPrice(dayConsumption domain.DayConsumption, region domain.Region) (domain.DayConsumptionWithPrice, error) {

	t, err := dayConsumption.GetTime()
	if err != nil {
		return domain.DayConsumptionWithPrice{}, err
	}

	dayPrices, err := ps.ElectricityPrice.GetPricePoints(t, region)
	if err != nil {
		return domain.DayConsumptionWithPrice{}, err
	}

	houreConsumptionsWithPrice := make([]domain.HourConsumptionWithPrice, 0)

	for i, hourConsumption := range dayConsumption.Consumption {
		houreConsumptionsWithPrice = append(houreConsumptionsWithPrice, domain.HourConsumptionWithPrice{
			Consumption: hourConsumption,
			Start:       dayPrices[i].TimeStart,
			End:         dayPrices[i].TimeEnd,
			SEKPerkWH:   dayPrices[i].SEKPerkWh,
			EURPerkWH:   dayPrices[i].EURPerkWh,
		})
	}

	return domain.DayConsumptionWithPrice{
		Date:            t,
		HourConsumption: houreConsumptionsWithPrice,
	}, nil
}
