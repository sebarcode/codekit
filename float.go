package codekit

import (
	"fmt"
	"math"
	"strings"
)

type FloatOpts struct {
	DecimalPoint int
	ThousandSep  string
	DecimalSep   string
	Rounding     string
}

func FormatFloat(num float64, opts *FloatOpts) string {
	if opts == nil {
		opts = &FloatOpts{
			DecimalPoint: 0,
			ThousandSep:  ",",
			DecimalSep:   ".",
			Rounding:     "",
		}
	}

	// apply rounding mode
	factor := math.Pow10(opts.DecimalPoint)
	switch strings.ToLower(opts.Rounding) {
	case "up":
		num = math.Ceil(num*factor) / factor
	case "down":
		num = math.Floor(num*factor) / factor
	default:
		num = math.Round(num*factor) / factor
	}

	neg := false
	if num < 0 {
		neg = true
		num = -num
	}

	// format with fixed decimals
	s := fmt.Sprintf("%.*f", opts.DecimalPoint, num)
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	// insert thousand separator
	var b strings.Builder
	n := len(intPart)
	for i := 0; i < n; i++ {
		// position from left
		if i > 0 && (n-i)%3 == 0 {
			b.WriteString(opts.ThousandSep)
		}
		b.WriteByte(intPart[i])
	}

	out := b.String()
	if opts.DecimalPoint > 0 {
		out = out + opts.DecimalSep + decPart
	}
	if neg {
		out = "-" + out
	}

	return out
}
