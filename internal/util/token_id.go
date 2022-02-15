package util

import (
	"fmt"
	"math/big"
	"strings"
)

func TokenIdToIntString(tokenId string) (string, error) {
	hexStr := strings.ReplaceAll(tokenId, "-", "")
	i := new(big.Int)
	_, success := i.SetString(hexStr, 16)
	if !success {
		return "", fmt.Errorf("failed to convert hex string to big int")
	}
	return i.String(), nil
}
