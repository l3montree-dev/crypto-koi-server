package server

import (
	"image"
	"image/draw"
	"path/filepath"
	"testing"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
	imageDraw "golang.org/x/image/draw"
)

func BenchmarkNearestNeighbor(b *testing.B) {
	basePath, err := filepath.Abs("../../images/raw")
	if err != nil {
		b.Fatal(err)
	}
	preloader := generator.NewMemoryPreloader(basePath)
	g := generator.NewGenerator(preloader)

	img, _ := g.TokenId2Image("243155001456846051896781451612598490801")

	scaledImg := image.NewRGBA(image.Rect(0, 0, 350, 350))

	for i := 0; i < b.N; i++ {
		imageDraw.NearestNeighbor.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)
	}
}

func BenchmarkBiLinear(b *testing.B) {
	basePath, err := filepath.Abs("../../images/raw")
	if err != nil {
		b.Fatal(err)
	}
	preloader := generator.NewMemoryPreloader(basePath)
	g := generator.NewGenerator(preloader)

	img, _ := g.TokenId2Image("243155001456846051896781451612598490801")

	scaledImg := image.NewRGBA(image.Rect(0, 0, 350, 350))

	for i := 0; i < b.N; i++ {
		imageDraw.BiLinear.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)
	}
}

func BenchmarkApproxBiLinear(b *testing.B) {
	basePath, err := filepath.Abs("../../images/raw")
	if err != nil {
		b.Fatal(err)
	}
	preloader := generator.NewMemoryPreloader(basePath)
	g := generator.NewGenerator(preloader)

	img, _ := g.TokenId2Image("243155001456846051896781451612598490801")

	scaledImg := image.NewRGBA(image.Rect(0, 0, 350, 350))

	for i := 0; i < b.N; i++ {
		imageDraw.ApproxBiLinear.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)
	}
}

func BenchmarkCatmullRom(b *testing.B) {
	basePath, err := filepath.Abs("../../images/raw")
	if err != nil {
		b.Fatal(err)
	}
	preloader := generator.NewMemoryPreloader(basePath)
	g := generator.NewGenerator(preloader)

	img, _ := g.TokenId2Image("243155001456846051896781451612598490801")

	scaledImg := image.NewRGBA(image.Rect(0, 0, 350, 350))

	for i := 0; i < b.N; i++ {
		imageDraw.CatmullRom.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)
	}
}
