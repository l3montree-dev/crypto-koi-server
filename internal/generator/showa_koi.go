package generator

import "image/color"

type ShowaKoi struct {
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

func (koi ShowaKoi) PrimaryColor() color.Color {
	return koi.redColor
}

func (koi ShowaKoi) GetFinImages(amount int, randomSeed int) []ImageWithColor {
	return pickAmount(
		amount,
		randomSeed,
		concatPreAllocate(
			withColor("fin", 1, 2, koi.redColor),
			withColor("fin", 1, 2, koi.blackColor),
		),
	)
}

func (koi ShowaKoi) AmountFinImages() (int, int) {
	return 0, 1
}

func (koi ShowaKoi) AmountHeadImages() (int, int) {
	return 0, 1
}

func (koi ShowaKoi) AmountBodyImages() (int, int) {
	return 1, 4
}

func (koi ShowaKoi) GetFinBackgroundColor(randomSeed int) color.Color {
	return koi.bodyAndFinColor
}

func (koi ShowaKoi) GetBodyColor(randomSeed int) color.Color {
	return koi.bodyAndFinColor
}

func (koi ShowaKoi) GetBodyImages(amount int, randomSeed int) []ImageWithColor {
	// generate the red color - so that all image patterns have the same red color
	return pickAmount(
		amount,
		randomSeed,
		concatPreAllocate(
			withColor("body", 1, 10, koi.redColor),
			withColor("body", 1, 10, koi.blackColor),
		),
	)
}

func (koi ShowaKoi) GetHeadImages(amount int, randomSeed int) []ImageWithColor {
	return pickAmount(amount, randomSeed, concatPreAllocate(
		withColor("head", 1, 5, koi.redColor),
		withColor("head", 1, 5, koi.redColor),
	))
}
