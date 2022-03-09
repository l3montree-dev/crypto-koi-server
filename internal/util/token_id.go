package util

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/google/uuid"
)

func UuidToUint256(tokenId string) (*big.Int, error) {
	hexStr := strings.ReplaceAll(tokenId, "-", "")

	uInt256, success := math.ParseBig256("0x" + hexStr)

	if !success {
		return nil, fmt.Errorf("failed to convert hex string to big int")
	}
	return uInt256, nil
}

func Uint256ToUuid(uInt *big.Int) (uuid.UUID, error) {
	hex := fmt.Sprintf("%x", uInt)
	for len(hex) < 32 {
		hex = "0" + hex
	}
	return uuid.Parse(hex)
}

func IsNotDigit(c rune) bool {
	return c < '0' || c > '9'
}

func IsHex(toCheck string) bool {
	for _, c := range toCheck {
		if IsNotDigit(c) {
			return false
		}
	}
	return true
}
