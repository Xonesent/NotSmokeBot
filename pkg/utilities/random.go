package utilities

import (
	"math/rand"
)

func RandomInRangeInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
