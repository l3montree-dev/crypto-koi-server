package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	stdMath "math"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/common/math"
)

const (
	MAX_BODY_PATTERNS int = 4
)

var (
	BOUNDS = image.Rect(0, 0, 1040, 1040)
)

type Generator struct {
	preloader Preloader
	koiTypes  []Koi
}

func NewGenerator(preloader Preloader) Generator {
	return Generator{
		preloader: preloader,
		koiTypes: []Koi{
			NewKohakuKoi(),
		},
	}
}

// will clone the image instead of mutating - this is important to keep the originals intact.
func applyColorToImage(c color.Color, img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, _, _, alphaChannel := img.At(x, y).RGBA()
			if alphaChannel != 0 {

				r, g, b, _ := c.RGBA()

				result.Set(x, y, color.RGBA{
					uint8(r), uint8(g), uint8(b), 255,
				})
			}
		}
	}
	return result
}

func (g *Generator) TokenId2Image(tokenId string) image.Image {
	// create a 39 long string - prepend with 0s if the provided tokenId is smaller
	necessaryPrefixes := 39 - len(tokenId)
	for i := 0; i < necessaryPrefixes; i++ {
		tokenId = "0" + tokenId
	}

	// convert the tokenId to a big integer
	tokenIdBigInt := math.MustParseBig256(tokenId)

	randomSource := rand.NewSource(tokenIdBigInt.Int64())
	r := rand.New(randomSource)
	// now the token id is 39 characters long.
	// extract all seed values. Just crop a few characters and convert them into integers.
	koiTypeSeed := r.Intn(len(g.koiTypes))
	bodyColorSeed := r.Intn(10)
	finColorSeed := r.Intn(255)
	amountOfBodyPatternsSeed := r.Intn(255)

	finPatternColorSeed := r.Intn(255)
	bodyPatternColorSeed := r.Intn(255)
	headPatternColorSeed := r.Intn(255)

	// start applying all seeds to first get the koy, and afterwards get all images.
	koi := g.koiTypes[koiTypeSeed%len(g.koiTypes)]
	finPatternSeed := r.Intn(len(koi.GetFinImages()))
	headPatternSeed := r.Intn(len(koi.GetHeadImages()))

	bodyColor := koi.GetBodyColorRange().Apply(bodyColorSeed)
	finColor := koi.GetFinBackgroundColorRange().Apply(finColorSeed)
	finPatternRaw := koi.GetFinImages()[finPatternSeed]
	headPatternRaw := koi.GetHeadImages()[headPatternSeed]

	// use a wait group to do stuff in parallel.
	wg := sync.WaitGroup{}
	wg.Add(4)
	staticImages := make(map[string]image.Image)

	go func() {
		defer wg.Done()
		staticImages["bodyImage"] = applyColorToImage(bodyColor, g.preloader.Body)
	}()
	go func() {
		defer wg.Done()
		staticImages["finImage"] = applyColorToImage(finColor, g.preloader.Fin)
	}()
	go func() {
		defer wg.Done()
		staticImages["finPattern"] = applyColorToImage(finPatternRaw.ColorRange.Apply(finPatternColorSeed), g.preloader.Images[finPatternRaw.ImageName])
	}()
	go func() {
		defer wg.Done()
		staticImages["headPattern"] = applyColorToImage(headPatternRaw.ColorRange.Apply(headPatternColorSeed), g.preloader.Images[headPatternRaw.ImageName])
	}()

	bodyPatterns := make([]image.Image, amountOfBodyPatternsSeed%int(stdMath.Min(float64(MAX_BODY_PATTERNS), float64(len(koi.GetBodyImages())))))

	for i := 0; i < len(bodyPatterns); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			img := koi.GetBodyImages()[(amountOfBodyPatternsSeed+i)%len(koi.GetBodyImages())]
			fmt.Println(amountOfBodyPatternsSeed+i, len(koi.GetBodyImages()), (amountOfBodyPatternsSeed+i)%len(koi.GetBodyImages()), img.ImageName)
			bodyPatterns[i] = applyColorToImage(img.ColorRange.Apply(bodyPatternColorSeed+i), g.preloader.Images[img.ImageName])
		}(i)
	}

	wg.Wait()

	// build the koi image.
	result := image.NewRGBA(BOUNDS)

	imageCollection := []image.Image{
		staticImages["bodyImage"],
		staticImages["finImage"],
		staticImages["finPattern"],
		staticImages["headPattern"],
	}

	imageCollection = append(imageCollection, bodyPatterns...)
	// the last image is the outline.
	imageCollection = append(imageCollection, g.preloader.Outline)

	// now we have all images in the collection.
	// we need to draw them in the correct order.
	for i := 0; i < len(imageCollection); i++ {
		img := imageCollection[i]
		bounds := img.Bounds()
		draw.Draw(result, bounds, img, bounds.Min, draw.Over)
	}

	return result
}
