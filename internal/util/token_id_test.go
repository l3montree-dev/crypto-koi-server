package util

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConvertCryptogotchi2TokenId(t *testing.T) {

	// the private key does not matter
	tokenId, err := UuidToUint256("b400af616cb4456589c4d6ba43f948b7")
	if err != nil {
		t.Fatal(err)
	}

	// "340282366920938463463374607431768211455"
	assert.Equal(t, "239264596381739575473221873891232270519", tokenId)
}

func TestSymmetry(t *testing.T) {
	expectedId, err := uuid.NewRandom()
	if err != nil {
		t.Fatal(err)
	}

	bigInt, err := UuidToUint256(expectedId.String())
	if err != nil {
		t.Fatal(err)
	}

	actualUuid, err := Uint256ToUuid(bigInt)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedId.String(), actualUuid.String())
}

func TestHexPrefixWith0(t *testing.T) {
	bigI := math.MustParseBig256("15077513574957955258113965027546720253")
	// will return the hex: b57d2d01d964fc3b2e3580a366567fd
	// which is invalid. To parse it, a leading 0 is required

	_, err := Uint256ToUuid(bigI)
	if err != nil {
		t.Fatal(err)
	}
}
