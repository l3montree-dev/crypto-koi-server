package generator

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"strconv"
	"sync"
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

type koiCtr = func(randomSeed int) Koi
type Generator struct {
	preloader Preloader
	koiCtrs   []koiCtr
}

func NewGenerator(preloader Preloader) Generator {
	return Generator{
		preloader: preloader,
		koiCtrs: []koiCtr{
			NewKohakuKoi,
			NewShowaKoi,
			NewUtsuriKoi,
			NewMonochromeKoi,
			NewShigureKoi,
		},
	}
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

func (g *Generator) GetKoi(tokenId string) (Koi, struct {
	r2 *rand.Rand
	r3 *rand.Rand
	r4 *rand.Rand
}) {
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
	koi := g.koiCtrs[r1.Intn(len(g.koiCtrs))](r1.Int())

	return koi, struct {
		r2 *rand.Rand
		r3 *rand.Rand
		r4 *rand.Rand
	}{r2, r3, r4}
}

func (g *Generator) TokenId2Image(tokenId string) (image.Image, Koi) {
	koi, randomizers := g.GetKoi(tokenId)
	r2, r3, r4 := randomizers.r2, randomizers.r3, randomizers.r4

	minBodyImages, maxBodyImages := koi.AmountBodyImages()
	minHeadImages, maxHeadImages := koi.AmountHeadImages()
	minFinImages, maxFinImages := koi.AmountFinImages()

	amountOfBodyImages := maxBodyImages
	if maxBodyImages != minBodyImages {
		// increment by 1 to include the max value into the possible values
		amountOfBodyImages = r2.Intn(maxBodyImages+1-minBodyImages) + minBodyImages
	}
	amountOfFinImages := maxFinImages
	if maxFinImages != minFinImages {
		// increment by 1 to include the max value into the possible values
		amountOfFinImages = r2.Intn(maxFinImages+1-minFinImages) + minFinImages
	}
	amountOfHeadImages := maxHeadImages
	if maxHeadImages != minHeadImages {
		// increment by 1 to include the max value into the possible values
		amountOfHeadImages = r2.Intn(maxHeadImages+1-minHeadImages) + minHeadImages
	}

	allImages := concatPreAllocate(
		// avoid providing zero to the r.Intn function. This will cause a panic.
		// therefore increase it to be at least 1 and then decrease it again after the call
		koi.GetBodyImages(amountOfBodyImages, r3.Intn(255)),
		koi.GetHeadImages(amountOfHeadImages, r3.Intn(255)),
		koi.GetFinImages(amountOfFinImages, r3.Intn(255)),
	)

	// add 2 for the body and fin image
	imgProcessingChan, imgResultChan := g.createChannels(len(allImages) + 2)

	g.startWorker(imgProcessingChan, imgResultChan)

	imgProcessingChan <- imageProcessingMessage{
		id:        0,
		baseImage: g.preloader.GetImage("body"),
		color:     koi.GetBodyColor(r4.Intn(255)),
	}

	imgProcessingChan <- imageProcessingMessage{
		id:        1,
		baseImage: g.preloader.GetImage("fins"),
		color:     koi.GetFinBackgroundColor(r4.Intn(255)),
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

	resultImages = append(resultImages, g.preloader.GetImage("outline"))
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
	if len(images) == 1 {
		// finished
		return images[0]
	}
	if len(images) == 2 {
		return combineImages(images[0], images[1])
	}

	return combineImages(images[0], recursiveBatchDraw(images[1:]))
}
