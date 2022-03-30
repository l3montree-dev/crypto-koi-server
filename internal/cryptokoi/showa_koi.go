package cryptokoi

import (
	"image/color"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type ShowaKoi struct {
	*CryptoKoi
	koiType KoiType
	// kohaku kois should have the same color for all patterns.
	redColor        color.Color
	bodyAndFinColor color.Color
	blackColor      color.Color
}

var _ Koi = ShowaKoi{}

func NewShowaKoi(randomSeed int) Koi {
	return ShowaKoi{
		koiType: Showa,
		// fix the pattern color.
		redColor:        RedColorRange.Apply(randomSeed),
		blackColor:      BlackColorRange.Apply(randomSeed),
		bodyAndFinColor: pickColorOutOf(randomSeed, RedColorRange, OrangeColorRange, WhiteColorRange),
	}
}

func (koi ShowaKoi) primaryColor() color.Color {
	return koi.redColor
}

func (koi ShowaKoi) getType() KoiType {
	return koi.koiType
}

func (koi ShowaKoi) getFinImages(amount int, randomSeed int) []util.ImageWithColor {
	return pickAmount(
		amount,
		randomSeed,
		util.ConcatPreAllocate(
			withColor("fin", 1, 2, koi.redColor),
			withColor("fin", 1, 2, koi.blackColor),
		),
	)
}

func (koi ShowaKoi) amountFinImages() (int, int) {
	return 0, 1
}

func (koi ShowaKoi) amountHeadImages() (int, int) {
	return 0, 1
}

func (koi ShowaKoi) amountBodyImages() (int, int) {
	return 1, 4
}

func (koi ShowaKoi) getFinBackgroundColor(randomSeed int) color.Color {
	return koi.bodyAndFinColor
}

func (koi ShowaKoi) getBodyColor(randomSeed int) color.Color {
	return koi.bodyAndFinColor
}

func (koi ShowaKoi) getBodyImages(amount int, randomSeed int) []util.ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return pickAmount(
		amount,
		randomSeed,
		util.ConcatPreAllocate(
			withColor("body", 1, 8, koi.redColor),
			withColor("body", 1, 8, koi.blackColor),
		),
	)
}

func (koi ShowaKoi) getHeadImages(amount int, randomSeed int) []util.ImageWithColor {
	return pickAmount(amount, randomSeed, util.ConcatPreAllocate(
		withColor("head", 1, 5, koi.redColor),
		withColor("head", 1, 5, koi.redColor),
	))
}
