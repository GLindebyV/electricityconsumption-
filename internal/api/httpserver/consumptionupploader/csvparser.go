package consumptionupploader

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/GLindebyV/electricityconsumption/internal/domain"
	"github.com/go-errors/errors"
)

const (
	cunsumptionPattern = `^.*(?P<date>\d{4}/\d{2}/\d{2}) \((?P<day>.{3,4})\)(?P<consumption>(,"\d{1,2},\d{1}"){24}),  \d*,  \d*,  (?P<total>\d*)`
)

var (
	errNoMatch         = errors.New("does not match pattern")
	errDataField       = errors.New("data field not found")
	errNotEnoughValues = errors.New("not enough values")
	errFloatConverting = errors.New("could not convert to float")
)

type dataRow struct {
	date             string
	day              string
	consumption      []string
	totalConsumption string
}

func (d *dataRow) toDayConsumption() (domain.DayConsumption, error) {
	consumption := make([]float64, 24)
	for i, c := range d.consumption {
		c = strings.Replace(c, ",", ".", 1)
		cf, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return domain.DayConsumption{}, fmt.Errorf("%w: %s", errFloatConverting, err)
		}
		consumption[i] = cf
	}

	totalConsumption, err := strconv.ParseFloat(d.totalConsumption, 64)
	if err != nil {
		return domain.DayConsumption{}, fmt.Errorf("%w: %s", errFloatConverting, err)
	}

	return domain.DayConsumption{
		Date:             d.date,
		Day:              d.day,
		Consumption:      consumption,
		TotalConsumption: totalConsumption,
	}, nil
}

func parseCunsumptionString(s string) (dataRow, error) {
	regExp := regexp.MustCompile(cunsumptionPattern)
	if !regExp.MatchString(s) {
		return dataRow{}, errNoMatch
	}

	fields := getRegexNamedGroupMap(regExp, s)

	date, ok := fields["date"]
	if !ok {
		return dataRow{}, errDataField
	}

	day, ok := fields["day"]
	if !ok {
		return dataRow{}, errDataField
	}

	consumption, ok := fields["consumption"]
	if !ok {
		return dataRow{}, errDataField
	}
	consumptions := strings.Split(consumption, `","`)
	for i, c := range consumptions {
		consumptions[i] = strings.Replace(c, `"`, "", -1)
	}

	if len(consumptions) != 24 {
		return dataRow{}, errNotEnoughValues
	}

	consumptions[0] = consumptions[0][1:]

	totalConsumption, ok := fields["total"]
	if !ok {
		return dataRow{}, errDataField
	}

	return dataRow{
		date:             date,
		day:              day,
		consumption:      consumptions,
		totalConsumption: totalConsumption,
	}, nil

}

func getRegexNamedGroupMap(exp *regexp.Regexp, str string) map[string]string {
	match := exp.FindStringSubmatch(str)

	paramsMap := map[string]string{}
	for i, name := range exp.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	return paramsMap
}

//regExp := regexp.MustCompile(recordConfig.Pattern)
//
//	if !regExp.MatchString(recordString) {
//		return nil, fmt.Errorf("row does not match %s", recordString)
//	}
//
//	fields := util.GetRegexNamedGroupMap(regExp, recordString)
//
//	for k, v := range fields {
//		fields[k] = strings.TrimRight(v, " ")
//	}
//
//	return fields, nil
