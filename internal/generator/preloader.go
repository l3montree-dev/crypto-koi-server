package generator

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// preload all images into ram for faster generation
type Preloader struct {
	Outline image.Image
	Body    image.Image
	Fin     image.Image
	Images  map[string]image.Image
}

func loadImage(basePath string, name string) image.Image {
	abs, err := filepath.Abs(filepath.Join(basePath, name))
	if err != nil {
		panic(err)
	}
	file, err := os.Open(abs)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

func NewPreloader(basePath string) Preloader {
	// load all images into ram
	preloader := Preloader{
		Outline: loadImage(basePath, "outline.png"),
		Fin:     loadImage(basePath, "fins.png"),
		Body:    loadImage(basePath, "body.png"),
		Images:  make(map[string]image.Image),
	}
	entries, err := os.ReadDir(basePath)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".png" {
			continue
		}
		preloader.Images[strings.TrimSuffix(entry.Name(), ".png")] = loadImage(basePath, entry.Name())
	}
	return preloader
}
