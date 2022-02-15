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

	for i := 0; i < 10; i++ {
		img := generator.TokenId2Image(fmt.Sprintf("%d", rand.Int()))
		file, _ := os.Create(fmt.Sprintf("empty_%d.png", i))
		png.Encode(file, img)
		file.Close()
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
	preloader := NewPreloader(path)
	generator := NewGenerator(preloader)

	for i := 0; i < b.N; i++ {
		generator.TokenId2Image(fmt.Sprintf("%d", rand.Int()))
	}
}
