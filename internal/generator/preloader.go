package generator

import (
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/image/draw"

	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

// preload all images into ram for faster generation
type MemoryPreloader struct {
	images map[string]image.Image
	// cached by size
	cache    sync.Map
	cacheMut sync.Mutex
	logger   *logrus.Entry
}

type Preloader interface {
	GetImage(imageName string, size int) image.Image
	BuildCachesForSizes(sizes []int) Preloader
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
		cache:  sync.Map{},
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

func (p *MemoryPreloader) BuildCachesForSizes(sizes []int) Preloader {
	wg := sync.WaitGroup{}
	wg.Add(len(sizes))
	for _, size := range sizes {
		go func(s int) {
			defer wg.Done()
			p.buildCacheForSize(s)
		}(size)
	}
	wg.Wait()
	// call the garbage collector - otherwise the makeTmpBuffer will consume lots of memory.
	runtime.GC()
	return p
}

func (p *MemoryPreloader) buildCacheForSize(size int) *MemoryPreloader {
	now := time.Now()
	var wg sync.WaitGroup

	cache, loaded := p.cache.LoadOrStore(size, &sync.Map{})
	if !loaded {
		// build the cache.
		for imageName := range p.images {
			wg.Add(1)
			go func(imgName string) {
				defer wg.Done()
				img := p.scaleImage(imgName, size)
				cache.(*sync.Map).Store(imgName, img)
			}(imageName)
		}
	}

	wg.Wait()
	p.logger.WithField("took", time.Since(now).String()).Info("cache built: ", size)
	return p
}

func (p *MemoryPreloader) scaleImage(imageName string, size int) image.Image {
	now := time.Now()
	// the image is not cached.
	rawImage, ok := p.images[imageName]
	if !ok {
		log.Fatalf("image: %s not preloaded", imageName)
	}

	// scale the image down.
	scaledImg := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(scaledImg, scaledImg.Rect, rawImage, rawImage.Bounds(), draw.Over, nil)
	p.logger.WithField("took", time.Since(now).String()).Warn("cache miss: scaled image: ", imageName, " to size: ", size)
	return scaledImg
}

func (p *MemoryPreloader) GetImage(imageName string, size int) image.Image {
	var img any
	var ok bool
	// check if the image is already cached
	if cache, loaded := p.cache.LoadOrStore(size, &sync.Map{}); loaded {
		if img, ok = cache.(*sync.Map).Load(imageName); ok {
			return img.(image.Image)
		}
		// the image does not exist in the size cache.
		// create it and cache it.
		img = p.scaleImage(imageName, size)
		cache.(*sync.Map).Store(imageName, img)
	} else {
		img = p.scaleImage(imageName, size)
		cache.(*sync.Map).Store(imageName, img)
	}
	return img.(image.Image)
}
