package generator

/*func TestT(t *testing.T) {
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

	t.Fail()
}
*/
