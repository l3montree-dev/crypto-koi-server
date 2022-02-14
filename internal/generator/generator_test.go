package generator

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestT(t *testing.T) {
	path, err := filepath.Abs(filepath.Join("..", "..", "images", "raw"))
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(42)
	preloader := NewPreloader(path)
	generator := NewGenerator(preloader)
	img := generator.TokenId2Image(fmt.Sprintf("%d", rand.Int()))
	file, _ := os.Create(fmt.Sprintf("empty_%d.png", 1))
	png.Encode(file, img)
	file.Close()
	t.Fail()
}
