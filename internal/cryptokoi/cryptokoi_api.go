package cryptokoi

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type CryptoKoiApi struct {
	privateKey *ecdsa.PrivateKey
	binding    *CryptoKoiBinding
}

func NewCryptokoiApi(privateHexKey string) CryptoKoiApi {
	privKey, err := crypto.HexToECDSA(strings.Replace(privateHexKey, "0x", "", 1))
	orchardclient.FailOnError(err, "Failed to parse private key")

	chainUrl := os.Getenv("CHAIN_URL")
	if chainUrl == "" {
		panic("CHAIN_URL is not set")
	}

	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	if contractAddress == "" {
		panic("CONTRACT_ADDRESS is not set")
	}

	client, err := ethclient.Dial(chainUrl)
	if err != nil {
		panic(err)
	}

	binding, err := NewCryptoKoiBinding(common.HexToAddress(contractAddress), client)
	if err != nil {
		panic(err)
	}

	return CryptoKoiApi{
		privateKey: privKey,
		binding:    binding,
	}
}

// pass the address in hex format
func (c *CryptoKoiApi) GetNftSignatureForCryptogotchi(cryptogotchi *models.Cryptogotchi, address string) (string, string, error) {
	tokenId, err := util.UuidToUint256(cryptogotchi.Id.String())
	if err != nil {
		return "", "", err
	}

	hash := solsha3.SoliditySHA3WithPrefix(solsha3.SoliditySHA3(
		// types
		[]string{"uint256", "address"},

		// values
		[]interface{}{
			tokenId,
			address,
		},
	))

	if err != nil {
		return "", "", err
	}

	sig, err := crypto.Sign(hash[:], c.privateKey)

	// this took ages: https://stackoverflow.com/questions/69762108/implementing-ethereum-personal-sign-eip-191-from-go-ethereum-gives-different-s
	// have a look at the link.
	sig[64] += 27

	if err != nil {
		return "", "", err
	}

	return hexutil.Encode(sig), tokenId.String(), nil
}

func (c *CryptoKoiApi) Redeem(address string, tokenId *big.Int, signature []byte) (*types.Transaction, error) {
	return c.binding.Redeem(&bind.TransactOpts{
		GasFeeCap: big.NewInt(10),
		GasTipCap: big.NewInt(10),
	}, common.HexToAddress(address), tokenId, signature)
}
