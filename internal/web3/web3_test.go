package web3

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
)

const (
	otherUserAddress  string = "0xa111C225A0aFd5aD64221B1bc1D5d817e5D3Ca15"
	privateKey        string = "0xc0c1e7d82fae79ce7727bd94e3e74deafbce52fc5618d9fd5557f41e83d4c149"
	expectedSignature string = "0x0577530589f065fdb25b8f29132865782ab2a4ea75a294ba56deecddeeefb77b18755f1811bb76dfadf417ff58f6bd2b593ddb4c80b1eaa85752e0df5a5b44f400"
)

func TestConvertCryptogotchi2TokenId(t *testing.T) {

	id := uuid.MustParse("b400af616cb4456589c4d6ba43f948b7")
	// the private key does not matter
	web3 := NewWeb3(privateKey)
	tokenId, err := web3.Uuid2Uint(id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "239264596381739575473221873891232270519", tokenId)
}

func TestGetSignature(t *testing.T) {
	web3 := NewWeb3(privateKey)
	cryptogotchi := &models.Cryptogotchi{
		Base: models.Base{Id: uuid.MustParse("b400af616cb4456589c4d6ba43f948b7")},
	}

	signature, _, err := web3.GetNftSignatureForCryptogotchi(cryptogotchi, otherUserAddress)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedSignature, signature)
}
