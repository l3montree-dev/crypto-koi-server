package cryptokoi

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
)

const (
	otherUserAddress  string = "0xa111C225A0aFd5aD64221B1bc1D5d817e5D3Ca15"
	privateKey        string = "0xc0c1e7d82fae79ce7727bd94e3e74deafbce52fc5618d9fd5557f41e83d4c149"
	expectedSignature string = "0x0577530589f065fdb25b8f29132865782ab2a4ea75a294ba56deecddeeefb77b18755f1811bb76dfadf417ff58f6bd2b593ddb4c80b1eaa85752e0df5a5b44f41b"
)

func TestRedeemToken(t *testing.T) {
	os.Setenv("CHAIN_URL", "http://localhost:8545")
	os.Setenv("CONTRACT_ADDRESS", "0x133c4b6c69322D09C5B266EFa9559173B6c9F029")
	cryptokoiApi := NewCryptokoiApi(privateKey)
	cryptogotchi := &models.Cryptogotchi{
		Base: models.Base{Id: uuid.MustParse("b400af616cb4456589c4d6ba43f948b7")},
	}

	signature, _, err := cryptokoiApi.GetNftSignatureForCryptogotchi(cryptogotchi, otherUserAddress)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedSignature, signature)
}
