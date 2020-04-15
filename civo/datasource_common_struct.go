package civo

// Customs struct for Filter used in dataSourceFiltersSchema
type Filter struct {
	Name   string
	Values []string
	Regex  bool
}
