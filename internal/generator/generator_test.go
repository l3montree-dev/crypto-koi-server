package generator

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
)

// BenchmarkGeneration-8   	       2	 609847462 ns/op	219332280 B/op	19943308 allocs/op
// BenchmarkGeneration-8   	       2	 640104734 ns/op	219332176 B/op	19943308 allocs/op
// BenchmarkGeneration-8   	       3	 460615302 ns/op	200263784 B/op	19942925 allocs/op
// changed the amount of images
// BenchmarkGeneration-8   	       7	 143402518 ns/op	71506013 B/op	 6510755 allocs/op
// BenchmarkGeneration-8   	       8	 138452379 ns/op	65436489 B/op	 6548596 allocs/op
// increasing the size of the images.
// BenchmarkGeneration-8   	       1	1388079368 ns/op	661901400 B/op	32137786 allocs/op

// BenchmarkGeneration-8   	       1	1562219833 ns/op	653115368 B/op	 1900504 allocs/op
// BenchmarkGeneration-8   	       1	3646803135 ns/op   1219412040 B/op	 1901191 allocs/op
func BenchmarkGeneration(b *testing.B) {
	path, _ := filepath.Abs(filepath.Join("..", "..", "images", "koi"))

	rand.Seed(42)
	preloader := NewMemoryPreloader(path).BuildCachesForSizes(500)
	generator := NewGenerator(preloader)

	for i := 0; i < b.N; i++ {
		generator.TokenId2Image(fmt.Sprintf("%d", rand.Int()), 500)
	}
}
