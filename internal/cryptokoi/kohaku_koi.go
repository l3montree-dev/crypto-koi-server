package cryptokoi

import (
	"image/color"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type KohakuKoi struct {
	*CryptoKoi
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

func (koi KohakuKoi) primaryColor() color.Color {
	return koi.color
}

func (koi KohakuKoi) getType() KoiType {
	return koi.koiType
}

func (koi KohakuKoi) getFinImages(amount int, randomSeed int) []util.ImageWithColor {
	return amountWithColor("fin", amount, randomSeed, 1, 2, koi.color)
}

func (koi KohakuKoi) amountFinImages() (int, int) {
	return 0, 1
}

func (koi KohakuKoi) amountHeadImages() (int, int) {
	return 0, 1
}

func (koi KohakuKoi) amountBodyImages() (int, int) {
	return 1, 4
}

func (koi KohakuKoi) getFinBackgroundColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi KohakuKoi) getBodyColor(randomSeed int) color.Color {
	return WhiteColorRange.Apply(randomSeed)
}

func (koi KohakuKoi) getBodyImages(amount int, randomSeed int) []util.ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return amountWithColor("body", amount, randomSeed, 1, 10, koi.color)
}

func (koi KohakuKoi) getHeadImages(amount int, randomSeed int) []util.ImageWithColor {
	return amountWithColor("head", amount, randomSeed, 1, 5, koi.color)
}
