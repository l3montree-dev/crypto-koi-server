package cryptokoi

import (
	"image/color"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type ShigureKoi struct {
	*CryptoKoi
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

func (koi ShigureKoi) primaryColor() color.Color {
	return koi.redColor
}

func (koi ShigureKoi) getType() KoiType {
	return koi.koiType
}

func (koi ShigureKoi) getFinImages(amount int, randomSeed int) []util.ImageWithColor {
	return amountWithColor("fin", amount, randomSeed, 1, 2, koi.blackColor)
}

func (koi ShigureKoi) amountFinImages() (int, int) {
	return 0, 1
}

func (koi ShigureKoi) amountHeadImages() (int, int) {
	return 1, 1
}

func (koi ShigureKoi) amountBodyImages() (int, int) {
	return 0, 2
}

func (koi ShigureKoi) getFinBackgroundColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi ShigureKoi) getBodyColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi ShigureKoi) getBodyImages(amount int, randomSeed int) []util.ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return amountWithColor("body", amount, randomSeed, 1, 8, koi.blackColor)
}

func (koi ShigureKoi) getHeadImages(amount int, randomSeed int) []util.ImageWithColor {
	return amountWithColor("head", amount, randomSeed, 6, 7, koi.redColor)
}
