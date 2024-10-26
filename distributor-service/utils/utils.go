package utils

import (
	"distributor-service/models"
	"strings"
)

func HasPermission(distributor *models.Distributor, regionCode string) bool {
	checkPermission := func(code string) bool {
		if permission, exists := distributor.Permissions[code]; exists {
			return permission
		}
		return false
	}

	if checkPermission(regionCode) {
		return true
	}

	parts := strings.Split(regionCode, "-")

	if len(parts) == 3 {
		province := parts[1] + "-" + parts[2]
		country := parts[2]
		return checkPermission(parts[0]+"-"+province) ||
			checkPermission(province) ||
			checkPermission(country)
	}

	if len(parts) == 2 {
		province := parts[0] + "-" + parts[1]
		country := parts[1]
		return checkPermission(province) ||
			checkPermission(country)
	}

	if len(parts) == 1 {
		return checkPermission(parts[0])
	}

	return false
}

func IsSubDistributor(to *models.Distributor, from *models.Distributor) bool {
	_, exists := from.SubDistributors[to.Name]
	return exists
}
