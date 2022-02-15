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

	Utsuri KoiType = "utsuri" // black pattern + white, yellow or red background
	// Asagi      KoiType = "asagi"      // blue colored scales red pattern
	Monochrome KoiType = "monochrome" // red, orange, yellow, yellow-greenish
	Shigure    KoiType = "shigure"    // white background + orange pattern
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
	OrangeColorRange ColorRange = ColorRange{
		raw: [3][2]int{
			{145, 255},
			{145, 255},
			{0, 75},
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
)

func pickColorOutOf(randomSeed int, ranges ...ColorRange) color.Color {
	r := rand.New(rand.NewSource(int64(randomSeed)))
	index := r.Intn(len(ranges))
	return ranges[index].Apply(randomSeed)
}
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
	PrimaryColor() color.Color
	GetType() KoiType
}

func pickAmount(amount, randomSeed int, images []ImageWithColor) []ImageWithColor {
	if amount == 0 {
		return []ImageWithColor{}
	}

	r := rand.New(rand.NewSource(int64(randomSeed)))

	result := make([]ImageWithColor, amount)
	for i := 0; i < amount && i < len(images); i++ {
		pickedIndex := r.Intn(len(images))
		result[i] = images[pickedIndex]
	}
	return result
}

func withColor(prefix string, possibleMin, possibleMax int, c color.Color) []ImageWithColor {
	// add +1 to include the possibleMax upper bound
	result := make([]ImageWithColor, (possibleMax-possibleMin)+1)
	idx := 0
	for i := possibleMin; i < possibleMax+1; i++ {
		result[idx] = ImageWithColor{
			ImageName: prefix + "_" + strconv.Itoa(i),
			Color:     c,
		}
		idx++
	}
	return result
}

func amountWithColor(prefix string, amount, randomSeed, possibleMin, possibleMax int, c color.Color) []ImageWithColor {
	if amount == 0 {
		return []ImageWithColor{}
	}

	possibilities := withColor(prefix, possibleMin, possibleMax, c)

	return pickAmount(amount, randomSeed, possibilities)
}
