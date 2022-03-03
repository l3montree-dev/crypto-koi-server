package util

import (
	"fmt"
	"image/color"
	"math"
)

func ConvertColor2Hex(c color.Color) string {
	r, g, b := RGBA(c)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func ConvertColor2HexWithoutHash(c color.Color) string {
	r, g, b := RGBA(c)
	return fmt.Sprintf("%02x%02x%02x", r, g, b)
}

func IsDark(c color.Color) bool {
	r, g, b := RGBA(c)
	hsp := math.Sqrt(0.299*float64(r*r) + 0.587*float64(g*g) + 0.114*float64(b*b))

	return hsp < 127.5
}

func RGBA(c color.Color) (uint32, uint32, uint32) {
	r, g, b, _ := c.RGBA()
	return r >> 8, g >> 8, b >> 8
}

func Shade(c color.Color, percentage int32) color.Color {
	r, g, b := RGBA(c)

	r = r * uint32(100+percentage) / 100
	g = g * uint32(100+percentage) / 100
	b = b * uint32(100+percentage) / 100

	r = uint32(math.Min(255, float64(r)))
	g = uint32(math.Min(255, float64(g)))
	b = uint32(math.Min(255, float64(b)))

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}
