package stores

import "distributor-service/models"

var Distributors map[string]*models.Distributor = make(map[string]*models.Distributor)
var CityToProvinceLookup map[string]string = nil
var ProvinceToCountryLookup map[string]string = nil
var RegionCodeToNameLookup map[string]string = nil

func GetDistributors() map[string]*models.Distributor {
	if Distributors == nil {
		return make(map[string]*models.Distributor)
	} else {
		return Distributors
	}

}

func GetDistributorByName(distributorName string) (*models.Distributor, bool) {
	if distributor, exists := Distributors[distributorName]; exists {
		return distributor, true
	} else {
		return nil, false
	}
}

func SaveDistributor(distributor *models.Distributor) {
	Distributors[distributor.Name] = distributor
}

func DeleteDistributor(distributorName string) {
	delete(Distributors, distributorName)
}

func DeleteSubDistributor(distributorName string, subDistributorName string) {
	delete(Distributors[distributorName].SubDistributors, subDistributorName)
	delete(Distributors[subDistributorName].ParentDistributors, distributorName)
}

func GetCityToProvinceLookup() map[string]string {
	if CityToProvinceLookup == nil {
		return make(map[string]string)
	}
	return CityToProvinceLookup
}

func SetCityToProvinceLookup(cityToProvinceCountry map[string]string) {
	CityToProvinceLookup = cityToProvinceCountry
}

func GetProvinceToCountryLookup() map[string]string {
	if ProvinceToCountryLookup == nil {
		return make(map[string]string)
	}
	return ProvinceToCountryLookup
}

func SetProvinceToCountryLookup(provinceToCountryLookup map[string]string) {
	ProvinceToCountryLookup = provinceToCountryLookup
}

func GetRegionCodeToNameLookup() map[string]string {
	if RegionCodeToNameLookup == nil {
		return make(map[string]string)
	}
	return RegionCodeToNameLookup
}

func SetRegionCodeToNameLookup(regionCodeToNameLookup map[string]string) {
	RegionCodeToNameLookup = regionCodeToNameLookup
}
