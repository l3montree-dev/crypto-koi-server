package web3

import (
	"crypto/ecdsa"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type Web3 struct {
	privateKey *ecdsa.PrivateKey
}

func NewWeb3(privateHexKey string) Web3 {
	privKey, err := crypto.HexToECDSA(strings.Replace(privateHexKey, "0x", "", 1))
	orchardclient.FailOnError(err, "Failed to parse private key")
	return Web3{
		privateKey: privKey,
	}
}

// pass the address in hex format
func (w *Web3) GetNftSignatureForCryptogotchi(cryptogotchi *models.Cryptogotchi, address string) (string, string, error) {
	tokenId, err := util.TokenIdToIntString(cryptogotchi.Id.String())
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

	sig, err := crypto.Sign(hash, w.privateKey)

	if err != nil {
		return "", "", err
	}

	return hexutil.Encode(sig), tokenId, nil
}
