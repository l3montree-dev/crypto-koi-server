package generator

import (
	"image/color"
)

type UtsuriKoi struct {
	koiType KoiType
	// Utsuri kois should have the same color for all patterns.
	color color.Color
}

var _ Koi = UtsuriKoi{}

func NewUtsuriKoi(randomSeed int) Koi {
	return UtsuriKoi{
		koiType: Utsuri,
		// fix the pattern color.
		color: BlackColorRange.Apply(randomSeed),
	}
}

func (koi UtsuriKoi) PrimaryColor() color.Color {
	return koi.color
}

func (koi UtsuriKoi) GetFinImages(amount int, randomSeed int) []ImageWithColor {
	return amountWithColor("fin", amount, randomSeed, 1, 2, koi.color)
}

func (koi UtsuriKoi) AmountFinImages() (int, int) {
	return 0, 1
}

func (koi UtsuriKoi) AmountHeadImages() (int, int) {
	return 0, 1
}

func (koi UtsuriKoi) AmountBodyImages() (int, int) {
	return 1, 4
}

func (koi UtsuriKoi) GetFinBackgroundColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi UtsuriKoi) GetBodyColor(randomSeed int) color.Color {
	return pickColorOutOf(randomSeed, WhiteColorRange, OrangeColorRange, RedColorRange)
}

func (koi UtsuriKoi) GetBodyImages(amount int, randomSeed int) []ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return amountWithColor("body", amount, randomSeed, 1, 10, koi.color)
}

func (koi UtsuriKoi) GetHeadImages(amount int, randomSeed int) []ImageWithColor {
	return amountWithColor("head", amount, randomSeed, 1, 5, koi.color)
}
