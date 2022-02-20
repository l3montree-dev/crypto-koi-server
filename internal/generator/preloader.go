package generator

import (
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

// preload all images into ram for faster generation
type MemoryPreloader struct {
	Images map[string]image.Image
}

type Preloader interface {
	GetImage(imageName string) image.Image
}

func loadImage(basePath string, name string) image.Image {
	abs, err := filepath.Abs(filepath.Join(basePath, name))
	if err != nil {
		orchardclient.Logger.Fatal(err)
	}
	file, err := os.Open(abs)
	if err != nil {
		orchardclient.Logger.Fatal(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		orchardclient.Logger.Fatal(err)
	}
	return img
}

func NewMemoryPreloader(basePath string) Preloader {
	// load all images into ram
	preloader := MemoryPreloader{
		Images: make(map[string]image.Image),
	}
	entries, err := os.ReadDir(basePath)
	if err != nil {
		orchardclient.Logger.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".png" {
			continue
		}
		preloader.Images[strings.TrimSuffix(entry.Name(), ".png")] = loadImage(basePath, entry.Name())
	}
	return &preloader
}

func (p *MemoryPreloader) GetImage(imageName string) image.Image {
	val, ok := p.Images[imageName]
	if !ok {
		log.Fatalf("image: %s not preloaded", imageName)
	}
	return val
}
