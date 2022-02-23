package cryptokoi

import (
	"image/color"
	"math/rand"
	"strconv"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
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

var EmptyImage = util.ImageWithColor{
	ImageName: "empty",
	Color:     color.White,
}

type KoiAttributes struct {
	KoiType      KoiType
	BodyImages   []util.ImageWithColor
	FinImages    []util.ImageWithColor
	HeadImages   []util.ImageWithColor
	BodyColor    color.Color
	FinColor     color.Color
	PrimaryColor color.Color
}
type Koi interface {
	getFinImages(amount int, randomSeed int) []util.ImageWithColor
	// [[r, r], [g, g], [b, b]]
	getFinBackgroundColor(randomSeed int) color.Color
	getBodyColor(randomSeed int) color.Color
	getBodyImages(amount int, randomSeed int) []util.ImageWithColor
	getHeadImages(amount int, randomSeed int) []util.ImageWithColor
	amountHeadImages() (int, int)
	amountBodyImages() (int, int)
	amountFinImages() (int, int)
	primaryColor() color.Color
	getType() KoiType
}

type koiCtr = func(randomSeed int) Koi

var koiCtrs = []koiCtr{
	NewKohakuKoi,
	NewShowaKoi,
	NewUtsuriKoi,
	NewMonochromeKoi,
	NewShigureKoi,
}

func NewKoi(tokenId string) CryptoKoi {
	// chunk the tokenId into 4 different sizes and create a random generator out of each.
	chunkSize := len(tokenId) / 4
	firstChunk, _ := strconv.ParseInt(tokenId[:chunkSize], 10, 64)
	secondChunk, _ := strconv.ParseInt(tokenId[chunkSize:chunkSize*2], 10, 64)
	thirdChunk, _ := strconv.ParseInt(tokenId[chunkSize*2:chunkSize*3], 10, 64)
	fourthChunk, _ := strconv.ParseInt(tokenId[chunkSize*3:], 10, 64)

	// this is just so random :-)
	r1, r2, r3, r4 := rand.New(rand.NewSource(firstChunk)), rand.New(rand.NewSource(secondChunk)), rand.New(rand.NewSource(thirdChunk)), rand.New(rand.NewSource(fourthChunk))
	// now the token id is 39 characters long.
	// extract all seed values. Just crop a few characters and convert them into integers.
	// start applying all seeds to first get the koy, and afterwards get all images.
	koi := koiCtrs[r1.Intn(len(koiCtrs))](r1.Int())

	return CryptoKoi{
		wrappedKoi: koi,
		randomizers: struct {
			r1 *rand.Rand
			r2 *rand.Rand
			r3 *rand.Rand
		}{
			// WE ALREADY USED THE FIRST RANDOMIZER to determine the type of the koi.
			r1: r2,
			r2: r3,
			r3: r4,
		},
	}
}

func pickAmount(amount, randomSeed int, images []util.ImageWithColor) []util.ImageWithColor {
	if amount == 0 {
		return []util.ImageWithColor{}
	}

	r := rand.New(rand.NewSource(int64(randomSeed)))

	result := make([]util.ImageWithColor, amount)
	for i := 0; i < amount && i < len(images); i++ {
		pickedIndex := r.Intn(len(images))
		result[i] = images[pickedIndex]
	}
	return result
}

func withColor(prefix string, possibleMin, possibleMax int, c color.Color) []util.ImageWithColor {
	// add +1 to include the possibleMax upper bound
	result := make([]util.ImageWithColor, (possibleMax-possibleMin)+1)
	idx := 0
	for i := possibleMin; i < possibleMax+1; i++ {
		result[idx] = util.ImageWithColor{
			ImageName: prefix + "_" + strconv.Itoa(i),
			Color:     c,
		}
		idx++
	}
	return result
}

func amountWithColor(prefix string, amount, randomSeed, possibleMin, possibleMax int, c color.Color) []util.ImageWithColor {
	if amount == 0 {
		return []util.ImageWithColor{}
	}

	possibilities := withColor(prefix, possibleMin, possibleMax, c)

	return pickAmount(amount, randomSeed, possibilities)
}
