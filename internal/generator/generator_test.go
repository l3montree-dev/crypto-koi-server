package generator

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneration(t *testing.T) {
	path, err := filepath.Abs(filepath.Join("..", "..", "images", "raw"))
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(42)
	preloader := NewMemoryPreloader(path)
	generator := NewGenerator(preloader)

	for i := 0; i < 10; i++ {
		f, _ := os.Create(fmt.Sprintf("%d.png", i))
		img := generator.TokenId2Image(fmt.Sprintf("%d", rand.Int()))
		png.Encode(f, img)
	}
	t.Fail()
}

func TestConsistency(t *testing.T) {
	path, err := filepath.Abs(filepath.Join("..", "..", "images", "raw"))
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(42)
	preloader := NewMemoryPreloader(path)
	generator := NewGenerator(preloader)

	hasher := sha256.New()

	img := generator.TokenId2Image(fmt.Sprintf("%d", 0))
	png.Encode(hasher, img)

	str := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	for i := 0; i < 5; i++ {
		hasher := sha256.New()

		img := generator.TokenId2Image(fmt.Sprintf("%d", 0))
		png.Encode(hasher, img)
		assert.Equal(t, str, base64.URLEncoding.EncodeToString(hasher.Sum(nil)))
	}
}

// BenchmarkGeneration-8   	       2	 609847462 ns/op	219332280 B/op	19943308 allocs/op
// BenchmarkGeneration-8   	       2	 640104734 ns/op	219332176 B/op	19943308 allocs/op
// BenchmarkGeneration-8   	       3	 460615302 ns/op	200263784 B/op	19942925 allocs/op
// changed the amount of images
// BenchmarkGeneration-8   	       7	 143402518 ns/op	71506013 B/op	 6510755 allocs/op
// BenchmarkGeneration-8   	       8	 138452379 ns/op	65436489 B/op	 6548596 allocs/op
func BenchmarkGeneration(b *testing.B) {
	path, _ := filepath.Abs(filepath.Join("..", "..", "images", "raw"))

	rand.Seed(42)
	preloader := NewMemoryPreloader(path)
	generator := NewGenerator(preloader)

	for i := 0; i < b.N; i++ {
		generator.TokenId2Image(fmt.Sprintf("%d", rand.Int()))
	}
}
