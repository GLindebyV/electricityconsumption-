package consumptionupploader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GLindebyV/electricityconsumption/internal/domain"
	"github.com/go-errors/errors"
)

type PriceService interface {
	GetDayConsumptionWithPrice(dayConsumption domain.DayConsumption, region domain.Region) (domain.DayConsumptionWithPrice, error)
}

type Controller struct {
	PriceService PriceService
}

func NewController(priceService PriceService) *Controller {
	return &Controller{
		PriceService: priceService,
	}
}

func (c *Controller) UpploadTest(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}

	fmt.Print(buf)

}

func (c *Controller) UpploadConsumption(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}

	dayConsumptions := make([]domain.DayConsumption, 0)

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		dr, err := parseCunsumptionString(line)
		if err != nil {
			if errors.Is(err, errNoMatch) {
				continue
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dc, err := dr.toDayConsumption()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dayConsumptions = append(dayConsumptions, dc)
	}

	dayConsumptionsWithPrice := make([]domain.DayConsumptionWithPrice, 0)

	for i, dc := range dayConsumptions {
		dcp, err := c.PriceService.GetDayConsumptionWithPrice(dc, domain.SE3)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Print(i)
			return
		}

		dayConsumptionsWithPrice = append(dayConsumptionsWithPrice, dcp)
	}

	resp := toResponse(dayConsumptionsWithPrice)

	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
