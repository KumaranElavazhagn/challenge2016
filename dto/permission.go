package dto

type CheckPermissionData struct {
	DistributorName string
	Regions         map[string]bool
}
