package generator

import (
	"image/color"
	"math/rand"
	"strconv"
)

type KoiType = string

const (
	Kohaku KoiType = "kohaku" // white background - red pattern
	Showa  KoiType = "showa"  // Kohaku but with black patterns

	Utsuri     KoiType = "utsuri"     // black pattern + white, yellow or red background
	Asagi      KoiType = "asagi"      // blue colored scales red pattern
	Monochrome KoiType = "monochrome" // red, orange, yellow, yellow-greenish
	Shigure    KoiType = "shigure"    // white background + orange pattern
	Sanke      KoiType = "sanke"      // white background + black pattern + red pattern on head only
)

// [[r, r], [g, g], [b, b]]
type ColorRange struct {
	raw [3][2]int
}

var (
	White ColorRange = ColorRange{
		raw: [3][2]int{
			{255, 255},
			{255, 255},
			{255, 255},
		},
	}
	RedColorRange ColorRange = ColorRange{
		raw: [3][2]int{
			{145, 255},
			{0, 75},
			{0, 54},
		},
	}
	WhiteColorRange ColorRange = ColorRange{
		raw: [3][2]int{
			{210, 255},
			{210, 255},
			{210, 255},
		},
	}
	BlackColorRange ColorRange = ColorRange{
		raw: [3][2]int{
			{0, 54},
			{0, 54},
			{0, 54},
		},
	}
	YellowColorRange ColorRange = ColorRange{
		raw: [3][2]int{
			{145, 255},
			{145, 255},
			{0, 75},
		},
	}
	OrangeColorRange ColorRange = ColorRange{
		raw: [3][2]int{
			{145, 255},
			{145, 255},
			{0, 75},
		},
	}
)

func limit(min, max, seed int) int {
	return min + seed%(max-min)
}

func (c ColorRange) Apply(randomSeed int) color.Color {
	r := limit(c.raw[0][0], c.raw[0][1], randomSeed)
	g := limit(c.raw[1][0], c.raw[1][1], randomSeed)
	b := limit(c.raw[2][0], c.raw[2][1], randomSeed)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

type ImageWithColor struct {
	// just the name of the image - without extension
	ImageName string
	// [[r, r], [g, g], [b, b]]
	// example: [[0, 255], [0, 255], [0, 255]]
	Color color.Color
}

var EmptyImage = ImageWithColor{
	ImageName: "empty",
	Color:     color.White,
}

type Koi interface {
	GetFinImages(amount int, randomSeed int) []ImageWithColor
	// [[r, r], [g, g], [b, b]]
	GetFinBackgroundColor(randomSeed int) color.Color
	GetBodyColor(randomSeed int) color.Color
	GetBodyImages(amount int, randomSeed int) []ImageWithColor
	GetHeadImages(amount int, randomSeed int) []ImageWithColor
	AmountHeadImages() (int, int)
	AmountBodyImages() (int, int)
	AmountFinImages() (int, int)
}

type KohakuKoi struct {
	koiType KoiType
}

var _ Koi = KohakuKoi{}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func amountWithColor(prefix string, amount, randomSeed, possibleMin, possibleMax int, c color.Color) []ImageWithColor {
	if amount == 0 {
		return []ImageWithColor{}
	}

	r := rand.New(rand.NewSource(int64(randomSeed)))

	possibleIndices := make([]int, possibleMax-possibleMin)
	for i := 0; i < len(possibleIndices); i++ {
		possibleIndices[i] = i + possibleMin
	}

	result := make([]ImageWithColor, amount)
	for i := 0; i < amount; i++ {
		pickedIndex := r.Intn(len(possibleIndices))
		result[i] = ImageWithColor{
			ImageName: prefix + "_" + strconv.Itoa(pickedIndex+1),
			Color:     c,
		}
		possibleIndices = remove(possibleIndices, pickedIndex)
	}
	return result
}

func NewKohakuKoi() Koi {
	return KohakuKoi{
		koiType: Kohaku,
	}
}

func (koi KohakuKoi) GetFinImages(amount int, randomSeed int) []ImageWithColor {
	res := []ImageWithColor{
		EmptyImage,
	}
	res = append(res, amountWithColor("fin", amount, randomSeed, 1, 2, RedColorRange.Apply(randomSeed))...)
	return res
}

func (koi KohakuKoi) AmountFinImages() (int, int) {
	return 0, 1
}

func (koi KohakuKoi) AmountHeadImages() (int, int) {
	return 1, 1
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
	return amountWithColor("body", amount, randomSeed, 1, 11, RedColorRange.Apply(randomSeed))
}

func (koi KohakuKoi) GetHeadImages(amount int, randomSeed int) []ImageWithColor {
	return amountWithColor("head", amount, randomSeed, 1, 5, RedColorRange.Apply(randomSeed))
}
