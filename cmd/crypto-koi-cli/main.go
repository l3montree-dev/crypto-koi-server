package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
)

// The cli can be used to generate koi images based upon the token id provided as the first argument.
func main() {
	err := godotenv.Load()

	drawPrimaryColor := flag.Bool("drawPrimaryColor", false, "draw the primary color onto the image")

	flag.Parse()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	baseImagePath := os.Getenv("BASE_IMAGE_PATH")

	if baseImagePath == "" {
		log.Fatal("BASE_IMAGE_PATH environment variable not set")
	}

	tokenId := flag.Arg(0)

	if tokenId == "" {
		log.Fatal("No token id provided as first argument (example: $ crypto-koi-cli <token id>)")
	}

	// generate the image based on the token id
	preloader := generator.NewMemoryPreloader(baseImagePath)
	g := generator.NewGenerator(preloader)

	img, koi := g.TokenId2Image(tokenId)

	if *drawPrimaryColor {
		primaryColor := koi.PrimaryColor()

		newImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
		for x := 0; x < newImg.Bounds().Max.X; x++ {
			for y := 0; y < newImg.Bounds().Max.Y; y++ {
				newImg.Set(x, y, primaryColor)
			}
		}

		draw.Draw(img.(draw.Image), img.Bounds(), newImg, image.Point{}, draw.Over)
	}

	f, err := os.Create(fmt.Sprintf("%s.png", tokenId))
	if err != nil {
		log.Fatal(err)
	}
	// encode the image to png
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}

	abs, err := filepath.Abs(fmt.Sprintf("%s.png", tokenId))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Image generated and saved at:", abs)
}
