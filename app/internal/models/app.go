package models

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gregaf/order-tracking-go/internal/util"
)

type FilterOperator string

var OperatorsOneValue = []FilterOperator{
	"eq",
	"ne",
	"gt",
	"lt",
}

var OperatorsTwoValues = []FilterOperator{
	"bt",
}

type FilterCriteria struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Values   []string       `json:"values"`
}

type GetResourceOptions struct {
	PageSize       int32  `json:"pageSize"`
	PageToken      string `json:"pageToken"`
	FilterCriteria []FilterCriteria
}

func (fc *FilterCriteria) Validate() error {

	if fc.Field == "" {
		return fmt.Errorf("field is required")
	}

	if fc.Operator == "" {
		return fmt.Errorf("operator is required")
	}

	if len(fc.Values) == 0 {
		return fmt.Errorf("values is required")
	}

	isOneValOp := util.IsOneOf(fc.Operator, OperatorsOneValue...)

	if isOneValOp && len(fc.Values) != 1 {
		return fmt.Errorf("operator %s can only have one value", fc.Operator)
	}

	isTwoValOp := util.IsOneOf(fc.Operator, OperatorsTwoValues...)

	if isTwoValOp && len(fc.Values) != 2 {
		return fmt.Errorf("operator %s can only have two values", fc.Operator)
	}

	if !isOneValOp && !isTwoValOp {
		return fmt.Errorf("operator %s is not supported", fc.Operator)
	}

	return nil
}

func ParseFilterCriteria(logger *slog.Logger, filterCriteria string) ([]FilterCriteria, error) {
	if filterCriteria == "" {
		return []FilterCriteria{}, nil
	}

	rawFilters := strings.Split(filterCriteria, ",")

	filters := []FilterCriteria{}
	for _, rawFilter := range rawFilters {
		parsedFilter := strings.Split(rawFilter, ":")

		// if parsedFilter has less than 3 elements, return error
		if len(parsedFilter) < 3 {
			return nil, fmt.Errorf("invalid filter criteria: %s", rawFilter)
		}

		filterField := strings.ToLower(parsedFilter[0])
		filterOperator := FilterOperator(parsedFilter[1])
		filterValues := parsedFilter[2:]

		filters = append(filters, FilterCriteria{
			Field:    filterField,
			Operator: filterOperator,
			Values:   filterValues,
		})
	}

	return filters, nil
}
