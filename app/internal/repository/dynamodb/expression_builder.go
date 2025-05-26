package dynamodbrepository

import (
	"fmt"

	"gregaf/order-tracking-go/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func buildFilterExpression(filters []models.FilterCriteria) (*expression.ConditionBuilder, error) {
	if len(filters) == 0 {
		return nil, nil
	}

	var filterExpr *expression.ConditionBuilder
	for _, f := range filters {

		newFilterExpr, err := getOperatorExpression(f)
		if err != nil {
			return nil, err
		}

		if filterExpr == nil {

			filterExpr = newFilterExpr
			continue
		}

		combinedFilterExpr := filterExpr.And(*newFilterExpr)
		filterExpr = &combinedFilterExpr
	}

	return filterExpr, nil
}

func getOperatorExpression(filter models.FilterCriteria) (*expression.ConditionBuilder, error) {

	switch filter.Operator {
	case "eq":
		expr := expression.Name(filter.Field).Equal(expression.Value(filter.Values[0]))
		return &expr, nil
	case "ne":
		expr := expression.Name(filter.Field).NotEqual(expression.Value(filter.Values[0]))
		return &expr, nil
	case "gt":
		expr := expression.Name(filter.Field).GreaterThan(expression.Value(filter.Values[0]))
		return &expr, nil
	case "lt":
		expr := expression.Name(filter.Field).LessThan(expression.Value(filter.Values[0]))
		return &expr, nil
	case "bt":
		expr := expression.Name(filter.Field).Between(expression.Value(filter.Values[0]), expression.Value(filter.Values[1]))
		return &expr, nil
	default:
		return nil, fmt.Errorf("operator %s is not supported", filter.Operator)
	}
}
