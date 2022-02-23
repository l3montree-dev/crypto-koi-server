package cryptokoi

import (
	"image/color"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type MonochromeKoi struct {
	*CryptoKoi
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

func (koi MonochromeKoi) primaryColor() color.Color {
	return koi.color
}

func (koi MonochromeKoi) getType() KoiType {
	return koi.koiType
}

func (koi MonochromeKoi) getFinImages(amount int, randomSeed int) []util.ImageWithColor {
	return []util.ImageWithColor{}
}

func (koi MonochromeKoi) amountFinImages() (int, int) {
	return 0, 0
}

func (koi MonochromeKoi) amountHeadImages() (int, int) {
	return 0, 0
}

func (koi MonochromeKoi) amountBodyImages() (int, int) {
	return 0, 0
}

func (koi MonochromeKoi) getFinBackgroundColor(randomSeed int) color.Color {
	return koi.color
}

func (koi MonochromeKoi) getBodyColor(randomSeed int) color.Color {
	return koi.color
}

func (koi MonochromeKoi) getBodyImages(amount int, randomSeed int) []util.ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return []util.ImageWithColor{}
}

func (koi MonochromeKoi) getHeadImages(amount int, randomSeed int) []util.ImageWithColor {
	return []util.ImageWithColor{}
}
