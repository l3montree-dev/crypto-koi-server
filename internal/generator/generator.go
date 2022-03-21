package generator

import (
	"image"
	"image/color"
	"image/draw"
	"sync"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/cryptokoi"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
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

type Generator struct {
	preloader Preloader
	debug     bool
}

func NewGenerator(preloader Preloader) Generator {
	return Generator{
		preloader: preloader,
	}
}

func (g *Generator) SetDebug(debug bool) {
	g.debug = debug
}

func (generator *Generator) createChannels(buffered int) (chan imageProcessingMessage, chan imageProcessingResult) {
	imageProcessingChan := make(chan imageProcessingMessage, buffered)
	imageProcessingResultChan := make(chan imageProcessingResult, buffered)

	return imageProcessingChan, imageProcessingResultChan
}
func (generator *Generator) startWorker(imageProcessingChan <-chan imageProcessingMessage, outputChan chan<- imageProcessingResult) {
	for i := 0; i < 1; i++ {
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
			if alphaChannel > 0 {
				result.Set(x, y, color.NRGBA{
					R: uint8((r) >> 8),
					G: uint8((g) >> 8),
					B: uint8((b) >> 8),
					A: uint8(alphaChannel >> 8),
				})
			}
		}
	}

	return result
}

func (g *Generator) TokenId2Image(tokenId string) (image.Image, *cryptokoi.CryptoKoi) {
	koi := cryptokoi.NewKoi(tokenId)

	attributes := koi.GetAttributes()
	allImages := util.ConcatPreAllocate(
		// avoid providing zero to the r.Intn function. This will cause a panic.
		// therefore increase it to be at least 1 and then decrease it again after the call
		attributes.BodyImages,
		attributes.HeadImages,
		attributes.FinImages,
	)
	// add 2 for the body and fin image
	imgProcessingChan, imgResultChan := g.createChannels(len(allImages) + 2)

	g.startWorker(imgProcessingChan, imgResultChan)

	imgProcessingChan <- imageProcessingMessage{
		id:        0,
		baseImage: g.preloader.GetImage("body"),
		color:     attributes.BodyColor,
	}

	imgProcessingChan <- imageProcessingMessage{
		id:        1,
		baseImage: g.preloader.GetImage("fins"),
		color:     attributes.FinColor,
	}

	for i, img := range allImages {
		imgProcessingChan <- imageProcessingMessage{
			id:        i + 2,
			baseImage: g.preloader.GetImage(img.ImageName),
			color:     img.Color,
		}
	}

	resultImages := make([]image.Image, len(allImages)+2)

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

	resultImages = append(resultImages, g.preloader.GetImage("highlights_1"), g.preloader.GetImage("outline"))
	// now we have all images in the collection.
	// we need to draw them in the correct order.
	result := recursiveBatchDraw(resultImages)
	return result, koi
}

func combineImages(dest image.Image, other image.Image) image.Image {
	draw.Draw(dest.(draw.Image), BOUNDS, other, BOUNDS.Min, draw.Over)
	return dest
}

func recursiveBatchDraw(images []image.Image) image.Image {
	result := image.NewRGBA(BOUNDS)
	for i := 0; i < len(images); i++ {
		combineImages(result, images[i])
	}
	return result
}
