package util

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"
)

func UuidToUint256(tokenId string) (*big.Int, error) {
	hexStr := strings.ReplaceAll(tokenId, "-", "")

	uInt256, success := math.ParseBig256("0x" + hexStr)

	if !success {
		return nil, fmt.Errorf("failed to convert hex string to big int")
	}
	return uInt256, nil
}
