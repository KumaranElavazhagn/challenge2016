package dto

type Distributor struct {
	Name    string
	Include map[string]bool
	Exclude map[string]bool
	Parent  string
}
