package generator

import (
	"image/color"
)

type KohakuKoi struct {
	koiType KoiType
	// kohaku kois should have the same color for all patterns.
	color color.Color
}

var _ Koi = KohakuKoi{}

func NewKohakuKoi(randomSeed int) Koi {
	return KohakuKoi{
		koiType: Kohaku,
		// fix the pattern color.
		color: RedColorRange.Apply(randomSeed),
	}
}

func (koi KohakuKoi) GetFinImages(amount int, randomSeed int) []ImageWithColor {
	return amountWithColor("fin", amount, randomSeed, 1, 2, koi.color)
}

func (koi KohakuKoi) AmountFinImages() (int, int) {
	return 0, 1
}

func (koi KohakuKoi) AmountHeadImages() (int, int) {
	return 0, 1
}

func (koi KohakuKoi) AmountBodyImages() (int, int) {
	return 1, 4
}

func (koi KohakuKoi) GetFinBackgroundColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi KohakuKoi) GetBodyColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi KohakuKoi) GetBodyImages(amount int, randomSeed int) []ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return amountWithColor("body", amount, randomSeed, 1, 10, koi.color)
}

func (koi KohakuKoi) GetHeadImages(amount int, randomSeed int) []ImageWithColor {
	return amountWithColor("head", amount, randomSeed, 1, 5, koi.color)
}
