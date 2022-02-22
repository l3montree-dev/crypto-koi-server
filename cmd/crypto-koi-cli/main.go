package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/db"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/service"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

func isNotDigit(c rune) bool {
	return c < '0' || c > '9'
}

func registerRandomUser(g *generator.Generator, amount int) {
	// register the user with the token id
	db, err := db.NewMySQL(db.MySQLConfig{
		User:     os.Getenv("DB_USER"),
		Password: strings.TrimSpace(string(util.MustReadFile(os.Getenv("DB_PASSWORD_FILE_PATH")))),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
	})

	if err != nil {
		log.Fatal(err)
	}

	userRep := repositories.NewGormUserRepository(db)
	refreshToken, _ := uuid.NewRandom()
	newUser := models.User{
		RefreshToken: refreshToken.String(),
	}
	userRep.Save(&newUser)

	cryptogotchiRep := repositories.NewGormCryptogotchiRepository(db)
	cryptogotchiSvc := service.NewCryptogotchiService(cryptogotchiRep, g)

	wg := sync.WaitGroup{}
	wg.Add(amount)
	for i := 0; i < amount; i++ {
		go func() {
			defer wg.Done()
			crypt, err := cryptogotchiSvc.GenerateCryptogotchiForUser(&newUser)
			if err != nil {
				log.Printf("WARNING - error occured: %e", err)
				return
			}
			log.Printf("User registered with id: %s and cryptogotchi tokenId: %s", newUser.Id, crypt.Id)
		}()
	}
	wg.Wait()
}

func drawImage(g *generator.Generator, drawPrimaryColor bool, tokenId string) {

	originalTokenId := tokenId
	if tokenId == "" {
		log.Println("No token id provided as first argument (example: $ crypto-koi-cli draw <token id>)")
		id, _ := uuid.NewRandom()
		tokenId = id.String()
		originalTokenId = tokenId
		log.Println("using random token id:", tokenId)
	}

	// check if the token id needs to be converted.
	if strings.IndexFunc(tokenId, isNotDigit) > -1 {
		// not only digits - use as hex.
		tmp, err := util.UuidToUint256(tokenId)
		if err != nil {
			log.Fatal(err)
		}

		tokenId = tmp.String()
	}

	img, koi := g.TokenId2Image(tokenId)

	if drawPrimaryColor {
		primaryColor := koi.PrimaryColor()

		newImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
		for x := 0; x < newImg.Bounds().Max.X; x++ {
			for y := 0; y < newImg.Bounds().Max.Y; y++ {
				newImg.Set(x, y, primaryColor)
			}
		}

		draw.Draw(img.(draw.Image), img.Bounds(), newImg, image.Point{}, draw.Over)
	}

	f, err := os.Create(fmt.Sprintf("%s.png", originalTokenId))
	if err != nil {
		log.Fatal(err)
	}
	// encode the image to png
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}

	abs, err := filepath.Abs(fmt.Sprintf("%s.png", originalTokenId))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Image generated and saved at:", abs)
}

// The cli can be used to generate koi images based upon the token id provided as the first argument.
func main() {

	err := godotenv.Load()

	baseImagePath := os.Getenv("BASE_IMAGE_PATH")

	if baseImagePath == "" {
		log.Fatal("BASE_IMAGE_PATH environment variable not set")
	}

	drawPrimaryColor := flag.Bool("drawPrimaryColor", false, "draw the primary color onto the image")
	debug := flag.Bool("debug", false, "enable debug mode")
	amount := flag.Int("amount", 1, "amount of users to register")

	flag.Parse()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// generate the image based on the token id
	preloader := generator.NewMemoryPreloader(baseImagePath)
	g := generator.NewGenerator(preloader)
	g.SetDebug(*debug)

	command := flag.Arg(0)

	switch command {
	case "draw":
		drawImage(&g, *drawPrimaryColor, flag.Arg(1))
	case "register":
		registerRandomUser(&g, *amount)
	case "sync-with-blockchain":
		// syncWithBlockchain()
	default:
		log.Fatalf("command: %s not found. Please use one of the following commands: register, draw", command)
	}
}
