package models

type DistributorNode struct {
	ParentDistributors map[string]bool
	SubDistributors    map[string]bool
}

type Distributor struct {
	Name        string          `json:"name"`
	Permissions map[string]bool `json:"permissions"`
	DistributorNode
}

type AddPermissionRequest struct {
	RegionCode string `json:"region_code"`
	IsIncluded bool   `json:"is_included"`
}

type DeletePermissionRequest struct {
	RegionCode string `json:"region_code"`
}

type CheckPermissionRequest struct {
	Name       string `json:"name"`
	RegionCode string `json:"region_code"`
}

type AuthorizeDistributorRequest struct {
	FromDistributor string `json:"from_distributor"`
	ToDistributor   string `json:"to_distributor"`
	RegionCode      string `json:"region_code"`
}

type SubDistributorRequest struct {
	DistributorName    string `json:"distributor_name"`
	SubDistributorName string `json:"sub_distributor_name"`
}
