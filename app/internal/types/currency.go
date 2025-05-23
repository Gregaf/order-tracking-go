package types

import "fmt"

type USD float64

func (c USD) String() string {
	return fmt.Sprintf("$ %.2f", c)
}
