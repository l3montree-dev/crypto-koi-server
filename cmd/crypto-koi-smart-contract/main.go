package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/cryptokoi"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

func main() {
	godotenv.Load()

	chainUrl := os.Getenv("CHAIN_URL")
	if chainUrl == "" {
		log.Fatal("CHAIN_URL environment variable is not defined")
	}
	chainWs := os.Getenv("CHAIN_WS")
	if chainWs == "" {
		log.Fatal("CHAIN_WS environment variable is not defined")
	}

	ethHttpClient, err := ethclient.Dial(chainUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer ethHttpClient.Close()

	ethWsClient, err := ethclient.Dial(chainWs)
	if err != nil {
		log.Fatal(err)
	}
	defer ethWsClient.Close()

	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	if contractAddress == "" {
		log.Fatal("CONTRACT_ADDRESS is not set")
	}

	httpBinding, err := cryptokoi.NewCryptoKoiBinding(common.HexToAddress(contractAddress), ethHttpClient)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("no token id specified")
	}

	tokenId := os.Args[1]

	if tokenId == "" {
		log.Fatal("tokenId is not set.")
	}

	if util.IsHex(tokenId) {
		uInt256, err := util.UuidToUint256(tokenId)
		if err != nil {
			log.Fatal(err)
		}
		tokenId = uInt256.String()
	}
	bigI, err := util.UuidToUint256(tokenId)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(httpBinding.TokenURI(nil, bigI))

}
