package generator

import (
	"image/color"
)

type MonochromeKoi struct {
	koiType KoiType
	// Monochrome kois should have the same color for all patterns.
	color color.Color
}

var _ Koi = MonochromeKoi{}

func NewMonochromeKoi(randomSeed int) Koi {
	return MonochromeKoi{
		koiType: Monochrome,
		// fix the pattern color.
		color: pickColorOutOf(randomSeed, WhiteColorRange, OrangeColorRange, RedColorRange, YellowColorRange),
	}
}

func (koi MonochromeKoi) PrimaryColor() color.Color {
	return koi.color
}

func (koi MonochromeKoi) GetFinImages(amount int, randomSeed int) []ImageWithColor {
	return []ImageWithColor{}
}

func (koi MonochromeKoi) AmountFinImages() (int, int) {
	return 0, 0
}

func (koi MonochromeKoi) AmountHeadImages() (int, int) {
	return 0, 0
}

func (koi MonochromeKoi) AmountBodyImages() (int, int) {
	return 0, 0
}

func (koi MonochromeKoi) GetFinBackgroundColor(randomSeed int) color.Color {
	return koi.color
}

func (koi MonochromeKoi) GetBodyColor(randomSeed int) color.Color {
	return koi.color
}

func (koi MonochromeKoi) GetBodyImages(amount int, randomSeed int) []ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return []ImageWithColor{}
}

func (koi MonochromeKoi) GetHeadImages(amount int, randomSeed int) []ImageWithColor {
	return []ImageWithColor{}
}
