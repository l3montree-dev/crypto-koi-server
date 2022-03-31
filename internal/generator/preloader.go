package generator

import (
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/image/draw"

	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

// preload all images into ram for faster generation
type MemoryPreloader struct {
	images map[string]image.Image
	// cached by size
	cache    map[int]map[string]image.Image
	cacheMut sync.Mutex
	logger   *logrus.Entry
}

type Preloader interface {
	GetImage(imageName string, size int) image.Image
	BuildCachesForSizes(sizes ...int) Preloader
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
		images: make(map[string]image.Image),
		cache:  make(map[int]map[string]image.Image),
		logger: orchardclient.Logger.WithField("component", "preloader"),
	}
	entries, err := os.ReadDir(basePath)
	if err != nil {
		orchardclient.Logger.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".png" {
			continue
		}
		preloader.images[strings.TrimSuffix(entry.Name(), ".png")] = loadImage(basePath, entry.Name())
	}
	return &preloader
}

func (p *MemoryPreloader) BuildCachesForSizes(sizes ...int) Preloader {
	wg := sync.WaitGroup{}
	wg.Add(len(sizes))
	for _, size := range sizes {
		go func(s int) {
			defer wg.Done()
			p.buildCacheForSize(s)
		}(size)
	}
	wg.Wait()
	return p
}

func (p *MemoryPreloader) buildCacheForSize(size int) *MemoryPreloader {
	if p.cache[size] == nil {
		p.cache[size] = make(map[string]image.Image)
	}
	for imageName := range p.images {
		img := p.scaleImage(imageName, size)
		p.cacheMut.Lock()
		p.cache[size][imageName] = img
		p.cacheMut.Unlock()
	}
	return p
}

func (p *MemoryPreloader) scaleImage(imageName string, size int) image.Image {
	p.logger.Warn("cache miss: scaling image: ", imageName, " to size: ", size)
	// the image is not cached.
	rawImage, ok := p.images[imageName]
	if !ok {
		log.Fatalf("image: %s not preloaded", imageName)
	}

	// scale the image down.
	scaledImg := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(scaledImg, scaledImg.Rect, rawImage, rawImage.Bounds(), draw.Over, nil)
	return scaledImg
}

func (p *MemoryPreloader) GetImage(imageName string, size int) image.Image {
	// check if the image is already cached
	if p.cache[size] != nil {
		if img, ok := p.cache[size][imageName]; ok {
			return img
		}
		// the image does not exist in the size cache.
		// create it and cache it.
		img := p.scaleImage(imageName, size)
		p.cacheMut.Lock()
		p.cache[size][imageName] = img
		p.cacheMut.Unlock()
	} else {
		img := p.scaleImage(imageName, size)
		p.cacheMut.Lock()
		// there is no cache for the size.
		// create one and cache it.
		p.cache[size] = make(map[string]image.Image)
		p.cache[size][imageName] = img
		p.cacheMut.Unlock()
	}
	return p.cache[size][imageName]
}
