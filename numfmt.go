package ddqp

import (
	"math"
	"strconv"
	"strings"
)

// formatFloatNoExp renders a float without scientific notation and trims
// unnecessary trailing zeros and decimal point.
func formatFloatNoExp(f float64) string {
	if f == math.Trunc(f) {
		return strconv.FormatFloat(f, 'f', 0, 64)
	}
	s := strconv.FormatFloat(f, 'f', 15, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}
