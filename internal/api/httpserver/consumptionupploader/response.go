package consumptionupploader

import "github.com/GLindebyV/electricityconsumption/internal/domain"

type DayConsumptionsWithPricingResponse struct {
	Data             []DayConsumptionWithPricing `json:"data"`
	TotalConsumption float64                     `json:"totalConsumption"`
	TotalPrice       float64                     `json:"totalPrice"`
}

type DayConsumptionWithPricing struct {
	Date                 string                     `json:"date"`
	ConsumptionWithPrice []HourConsumptionWithPrice `json:"consumptionWithPrice"`
	TotalConsumption     float64                    `json:"totalConsumption"`
	TotalPrice           float64                    `json:"totalPrice"`
}

type HourConsumptionWithPrice struct {
	Consumption float64 `json:"consumption"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
	SEKPerkWH   float64 `json:"SEKPerkWH"`
	Price       float64 `json:"price"`
}

func toResponse(dayConsumptionsWithPrice []domain.DayConsumptionWithPrice) DayConsumptionsWithPricingResponse {
	response := DayConsumptionsWithPricingResponse{
		Data: make([]DayConsumptionWithPricing, 0),
	}

	totalConsumption := 0.0
	totalPrice := 0.0

	for _, dayConsumptionWithPrice := range dayConsumptionsWithPrice {
		response.Data = append(response.Data, toDayConsumptionWithPricing(dayConsumptionWithPrice))
		totalConsumption += dayConsumptionWithPrice.GetTotalConsumption()
		totalPrice += dayConsumptionWithPrice.GetTotalPriceSEK()
	}

	response.TotalConsumption = totalConsumption
	response.TotalPrice = totalPrice

	return response
}

func toDayConsumptionWithPricing(dayConsumptionWithPrice domain.DayConsumptionWithPrice) DayConsumptionWithPricing {
	consumptionWithPrice := make([]HourConsumptionWithPrice, 0)

	for _, hourConsumptionWithPrice := range dayConsumptionWithPrice.HourConsumption {
		consumptionWithPrice = append(consumptionWithPrice, toHourConsumptionWithPrice(hourConsumptionWithPrice))
	}

	return DayConsumptionWithPricing{
		Date:                 dayConsumptionWithPrice.Date.Format("2006-01-02"),
		ConsumptionWithPrice: consumptionWithPrice,
		TotalConsumption:     dayConsumptionWithPrice.GetTotalConsumption(),
		TotalPrice:           dayConsumptionWithPrice.GetTotalPriceSEK(),
	}
}

func toHourConsumptionWithPrice(hourConsumptionWithPrice domain.HourConsumptionWithPrice) HourConsumptionWithPrice {
	return HourConsumptionWithPrice{
		Consumption: hourConsumptionWithPrice.Consumption,
		Start:       hourConsumptionWithPrice.Start.Format("15:04"),
		End:         hourConsumptionWithPrice.End.Format("15:04"),
		SEKPerkWH:   hourConsumptionWithPrice.SEKPerkWH,
		Price:       hourConsumptionWithPrice.GetTotalPriceSEK(),
	}
}
