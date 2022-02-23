package cryptokoi

import (
	"crypto/ecdsa"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type CryptoKoiApi struct {
	privateKey *ecdsa.PrivateKey
	binding    *CryptoKoiBinding
	logger     *logrus.Entry
}

func NewCryptokoiApi(privateHexKey string, binding *CryptoKoiBinding) CryptoKoiApi {
	privKey, err := crypto.HexToECDSA(strings.Replace(privateHexKey, "0x", "", 1))
	logger := orchardclient.Logger.WithField("component", "CryptoKoiApi")
	if err != nil {
		logger.Fatal(err)
	}

	return CryptoKoiApi{
		privateKey: privKey,
		binding:    binding,
		logger:     logger,
	}
}

// pass the address in hex format
func (c *CryptoKoiApi) GetNftSignatureForCryptogotchi(cryptogotchiId string, address string) (string, string, error) {
	tokenId, err := util.UuidToUint256(cryptogotchiId)
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
