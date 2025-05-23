package types

import "fmt"

type Weight struct {
	Value float64 `json:"value" dynamodbav:"Value"`
	Unit  string  `json:"unit" dynamodbav:"Unit"`
}

const (
	KG_PER_LB float64 = 0.45359237
	LB_PER_KG float64 = 2.20462
)

func (w Weight) ToKilograms() Weight {
	if w.Unit == "lbs" {
		return Weight{Value: w.Value * KG_PER_LB, Unit: "kg"}
	}

	return w
}

func (w Weight) ToPounds() Weight {
	if w.Unit == "kg" {
		return Weight{Value: w.Value * LB_PER_KG, Unit: "lbs"}
	}

	return w
}

func (w Weight) String() string {
	return fmt.Sprintf("%.2f %s", w.Value, w.Unit)
}
