package utils

import (
	"encoding/csv"
	"os"

	"distributor-service/stores"
)

func LoadCityData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	cityToProvinceLookup := stores.GetCityToProvinceLookup()
	provinceToCountryLookup := stores.GetProvinceToCountryLookup()
	regionCodeToNameLookup := stores.GetRegionCodeToNameLookup()
	for _, record := range records {
		cityCode := record[0] + "-" + record[1] + "-" + record[2]
		provinceCode := record[1] + "-" + record[2]
		countryCode := record[2]
		cityToProvinceLookup[cityCode] = provinceCode
		provinceToCountryLookup[provinceCode] = countryCode
		regionCodeToNameLookup[cityCode] = cityCode
		regionCodeToNameLookup[provinceCode] = provinceCode
		regionCodeToNameLookup[countryCode] = countryCode
	}
	stores.SetCityToProvinceLookup(cityToProvinceLookup)
	stores.SetProvinceToCountryLookup(provinceToCountryLookup)
	stores.SetRegionCodeToNameLookup(regionCodeToNameLookup)

	return nil
}
