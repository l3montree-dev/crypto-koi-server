package generator

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/math"
)

const (
	MAX_BODY_PATTERNS int = 4
)

var (
	BOUNDS = image.Rect(0, 0, 1040, 1040)
)

type imageProcessingMessage struct {
	baseImage image.Image
	color     color.Color
	id        int
}

type imageProcessingResult struct {
	result image.Image
	id     int
}

func concatPreAllocate(slices ...[]ImageWithColor) []ImageWithColor {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]ImageWithColor, totalLen)

	i := 0
	for _, s := range slices {
		for _, img := range s {
			tmp[i] = img
			i++
		}
	}
	return tmp
}

type koiCtr = func() Koi
type Generator struct {
	preloader Preloader
	koiCtrs   []koiCtr
}

func NewGenerator(preloader Preloader) Generator {
	return Generator{
		preloader: preloader,
		koiCtrs: []koiCtr{
			NewKohakuKoi,
		},
	}
}

func (generator *Generator) createChannels(buffered int) (chan imageProcessingMessage, chan imageProcessingResult) {
	imageProcessingChan := make(chan imageProcessingMessage, buffered)
	imageProcessingResultChan := make(chan imageProcessingResult, buffered)

	return imageProcessingChan, imageProcessingResultChan
}
func (generator *Generator) startWorker(imageProcessingChan <-chan imageProcessingMessage, outputChan chan<- imageProcessingResult) {
	for i := 0; i < 10; i++ {
		go func() {
			for msg := range imageProcessingChan {
				outputChan <- imageProcessingResult{
					result: generator.applyColorToImage(msg.color, msg.baseImage),
					id:     msg.id,
				}
			}
		}()
	}
}

// will clone the image instead of mutating - this is important to keep the originals intact.
func (generator *Generator) applyColorToImage(c color.Color, img image.Image) image.Image {
	bounds := img.Bounds()
	r, g, b, _ := c.RGBA()
	result := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, _, _, alphaChannel := img.At(x, y).RGBA()
			if alphaChannel>>8 > 240 {
				result.Set(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(alphaChannel >> 8),
				})
			}

		}

	}

	return result
}

func (g *Generator) TokenId2Image(tokenId string) image.Image {
	// convert the tokenId to a big integer
	tokenIdBigInt := math.MustParseBig256(tokenId)
	randomSource := rand.NewSource(tokenIdBigInt.Int64())
	r := rand.New(randomSource)
	// now the token id is 39 characters long.
	// extract all seed values. Just crop a few characters and convert them into integers.
	// start applying all seeds to first get the koy, and afterwards get all images.
	koi := g.koiCtrs[r.Intn(len(g.koiCtrs))]()

	minBodyImages, maxBodyImages := koi.AmountBodyImages()
	minHeadImages, maxHeadImages := koi.AmountHeadImages()
	minFinImages, maxFinImages := koi.AmountFinImages()

	amountOfBodyImages := maxBodyImages
	if maxBodyImages != minBodyImages {
		amountOfBodyImages = r.Intn(maxBodyImages-minBodyImages) + minBodyImages
	}
	amountOfFinImages := maxFinImages
	if maxFinImages != minFinImages {
		amountOfFinImages = r.Intn(maxFinImages-minFinImages) + minFinImages
	}
	amountOfHeadImages := maxHeadImages
	if maxHeadImages != minHeadImages {
		amountOfHeadImages = r.Intn(maxHeadImages-minHeadImages) + minHeadImages
	}

	allImages := concatPreAllocate(
		// avoid providing zero to the r.Intn function. This will cause a panic.
		// therefore increase it to be at least 1 and then decrease it again after the call
		koi.GetBodyImages(amountOfBodyImages, r.Intn(255)),
		koi.GetHeadImages(amountOfHeadImages, r.Intn(255)),
		koi.GetFinImages(amountOfFinImages, r.Intn(255)),
	)

	// add 2 for the body and fin image
	imgProcessingChan, imgResultChan := g.createChannels(len(allImages) + 2)

	g.startWorker(imgProcessingChan, imgResultChan)

	imgProcessingChan <- imageProcessingMessage{
		id:        0,
		baseImage: g.preloader.GetImage("body"),
		color:     koi.GetBodyColor(r.Intn(255)),
	}

	imgProcessingChan <- imageProcessingMessage{
		id:        1,
		baseImage: g.preloader.GetImage("fins"),
		color:     koi.GetBodyColor(r.Intn(255)),
	}

	for i, img := range allImages {
		imgProcessingChan <- imageProcessingMessage{
			id:        i + 2,
			baseImage: g.preloader.GetImage(img.ImageName),
			color:     img.Color,
		}
	}

	resultImages := make([]image.Image, len(allImages)+2)

	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(len(resultImages))
	// collect all images again.
	go func() {
		for img := range imgResultChan {
			resultImages[img.id] = img.result
			wg.Done()
		}
	}()

	wg.Wait()
	close(imgProcessingChan)
	close(imgResultChan)

	resultImages = append(resultImages, g.preloader.GetImage("outline"))
	elapsed := time.Since(start)
	log.Printf("collecting took %s", elapsed)
	start = time.Now()
	// now we have all images in the collection.
	// we need to draw them in the correct order.

	result := recursiveBatchDraw(resultImages)

	elapsed = time.Since(start)
	log.Printf("Drawing took %s", elapsed)
	return result
}

func combineImages(dest image.Image, other image.Image) image.Image {
	draw.Draw(dest.(draw.Image), BOUNDS, other, BOUNDS.Min, draw.Over)
	return dest
}

func recursiveBatchDraw(images []image.Image) image.Image {
	if len(images) == 1 {
		// finished
		return images[0]
	}
	if len(images) == 2 {
		return combineImages(images[0], images[1])
	}

	return combineImages(images[0], recursiveBatchDraw(images[1:]))
}
