package cryptokoi

import (
	"image/color"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type UtsuriKoi struct {
	*CryptoKoi
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

func (koi UtsuriKoi) primaryColor() color.Color {
	return koi.color
}

func (koi UtsuriKoi) getType() KoiType {
	return koi.koiType
}

func (koi UtsuriKoi) getFinImages(amount int, randomSeed int) []util.ImageWithColor {
	return amountWithColor("fin", amount, randomSeed, 1, 2, koi.color)
}

func (koi UtsuriKoi) amountFinImages() (int, int) {
	return 0, 1
}

func (koi UtsuriKoi) amountHeadImages() (int, int) {
	return 0, 1
}

func (koi UtsuriKoi) amountBodyImages() (int, int) {
	return 1, 4
}

func (koi UtsuriKoi) getFinBackgroundColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi UtsuriKoi) getBodyColor(randomSeed int) color.Color {
	return pickColorOutOf(randomSeed, WhiteColorRange, OrangeColorRange, RedColorRange)
}

func (koi UtsuriKoi) getBodyImages(amount int, randomSeed int) []util.ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return amountWithColor("body", amount, randomSeed, 1, 8, koi.color)
}

func (koi UtsuriKoi) getHeadImages(amount int, randomSeed int) []util.ImageWithColor {
	return amountWithColor("head", amount, randomSeed, 1, 5, koi.color)
}
