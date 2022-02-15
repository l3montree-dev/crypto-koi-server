package generator

import (
	"image/color"
)

type ShigureKoi struct {
	koiType    KoiType
	blackColor color.Color
	redColor   color.Color
}

var _ Koi = ShigureKoi{}

func NewShigureKoi(randomSeed int) Koi {
	return ShigureKoi{
		koiType: Shigure,
		// fix the pattern color.
		blackColor: BlackColorRange.Apply(randomSeed),
		redColor:   RedColorRange.Apply(randomSeed),
	}
}

func (koi ShigureKoi) PrimaryColor() color.Color {
	return koi.redColor
}

func (koi ShigureKoi) GetType() KoiType {
	return koi.koiType
}

func (koi ShigureKoi) GetFinImages(amount int, randomSeed int) []ImageWithColor {
	return amountWithColor("fin", amount, randomSeed, 1, 2, koi.blackColor)
}

func (koi ShigureKoi) AmountFinImages() (int, int) {
	return 0, 1
}

func (koi ShigureKoi) AmountHeadImages() (int, int) {
	return 1, 1
}

func (koi ShigureKoi) AmountBodyImages() (int, int) {
	return 0, 2
}

func (koi ShigureKoi) GetFinBackgroundColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi ShigureKoi) GetBodyColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi ShigureKoi) GetBodyImages(amount int, randomSeed int) []ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return amountWithColor("body", amount, randomSeed, 1, 10, koi.blackColor)
}

func (koi ShigureKoi) GetHeadImages(amount int, randomSeed int) []ImageWithColor {
	return amountWithColor("head", amount, randomSeed, 6, 7, koi.redColor)
}
