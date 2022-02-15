package util

import (
	"fmt"
	"image/color"
)

func ConvertColor2Hex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
