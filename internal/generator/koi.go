package generator

import (
	"image/color"
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
			{190, 255},
			{190, 255},
			{190, 255},
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

type ImageWithColorRange struct {
	// just the name of the image - without extension
	ImageName string
	// [[r, r], [g, g], [b, b]]
	// example: [[0, 255], [0, 255], [0, 255]]
	ColorRange ColorRange
}

var EmptyImage = ImageWithColorRange{
	ImageName:  "empty",
	ColorRange: WhiteColorRange,
}

type Koi interface {
	GetFinImages() []ImageWithColorRange
	// [[r, r], [g, g], [b, b]]
	GetFinBackgroundColorRange() ColorRange
	GetBodyColorRange() ColorRange
	GetBodyImages() []ImageWithColorRange
	GetHeadImages() []ImageWithColorRange
}

type KohakuKoi struct {
	koiType KoiType
}

func allWithColorRange(prefix string, min, max int, colorRange ColorRange) []ImageWithColorRange {
	result := make([]ImageWithColorRange, max-min)
	for i := 0; i < max-min; i++ {
		result[i] = ImageWithColorRange{
			ImageName:  prefix + "_" + strconv.Itoa(min+i),
			ColorRange: colorRange,
		}
	}
	return result
}

func NewKohakuKoi() KohakuKoi {
	return KohakuKoi{
		koiType: Kohaku,
	}
}

func (koi KohakuKoi) GetFinImages() []ImageWithColorRange {
	res := []ImageWithColorRange{
		EmptyImage,
	}
	res = append(res, allWithColorRange("fin", 1, 2, RedColorRange)...)
	return res
}

func (koi KohakuKoi) GetFinBackgroundColorRange() ColorRange {
	return WhiteColorRange
}

func (koi KohakuKoi) GetBodyColorRange() ColorRange {
	return WhiteColorRange
}

func (koi KohakuKoi) GetBodyImages() []ImageWithColorRange {
	res := []ImageWithColorRange{}
	res = append(res, allWithColorRange("body", 1, 11, RedColorRange)...)
	return res
}

func (koi KohakuKoi) GetHeadImages() []ImageWithColorRange {
	res := []ImageWithColorRange{}
	res = append(res, allWithColorRange("head", 1, 5, RedColorRange)...)
	return res
}
