package util

import "image/color"

type ImageWithColor struct {
	// just the name of the image - without extension
	ImageName string
	// [[r, r], [g, g], [b, b]]
	// example: [[0, 255], [0, 255], [0, 255]]
	Color color.Color
}

func ConcatPreAllocate(slices ...[]ImageWithColor) []ImageWithColor {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]ImageWithColor, totalLen)

	i := 0
	for _, s := range slices {
		for _, img := range s {
			tmp[i] = img
			i++
		}
	}
	return tmp
}
