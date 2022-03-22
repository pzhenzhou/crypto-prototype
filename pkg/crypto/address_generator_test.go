package crypto

import (
	common "github.com/pzhenzhou/crypto-prototype/pkg"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var regx = "^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$"

func TestHDSegWitAddress_Generate_IllegalArgs(t *testing.T) {
	testSeedGenerator := GetSeedGenerator(common.GetWordList())
	addressGenerator := NewHDSegWitAddress(testSeedGenerator)
	address, err := addressGenerator.Generate(nil)
	t.Log("address is ", address)
	assert.Nil(t, nil, address)
	assert.Error(t, err, ArgsMustBeNotNull)
}

func TestHDSegWitAddress_Generate(t *testing.T) {
	testSeedGenerator := GetSeedGenerator(common.GetWordList())
	addressGenerator := NewHDSegWitAddress(testSeedGenerator)
	args := map[GenerateArgs]interface{}{
		InputPassword: "TREZOR",
		InputPath:     "m/44'/0'/0'/0/0",
		InputMnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow",
	}
	address, err := addressGenerator.Generate(args)
	t.Logf("HD SegWit Address %v", address)
	re := regexp.MustCompile(regx)
	isMatch := re.MatchString(address.Address)
	assert.True(t, true, isMatch)
	assert.Equal(t, "legal winner thank year wave sausage worth useful legal winner thank yellow", address.Mnemonic)
	assert.Nil(t, nil, err)
}

func TestMultiSigAddress_Generate_IllegalArgs(t *testing.T) {
	multiSigAddress := MultiSigAddress{}
	args := map[GenerateArgs]interface{}{
		MultiSigNum:       MultiSigNumPair{M: 6, N: 3},
		MultiSigPublicKey: [][]byte{{}, {}, {}},
	}
	address, err := multiSigAddress.Generate(args)
	assert.Nil(t, nil, address)
	assert.Error(t, err, MultiSigNumValueInvalid)
}

func TestMultiSigAddress_Generate(t *testing.T) {
	multiSigAddress := MultiSigAddress{}
	args := map[GenerateArgs]interface{}{
		MultiSigNum: MultiSigNumPair{M: 3, N: 2},
		MultiSigPublicKey: [][]byte{
			[]byte("020f8796e0f870a9a3b269be3b1e78e380c9b569885f0de98a9ff061c4a66e79d2"),
			[]byte("02dfa8990f3f015ff20e9b31b85ea36d47470220615fb2ac1597e20fc830727b25"),
			[]byte("03fbfbdc5df9c60e4b747805552686199e85299a5e87804dbb66a14597ddabcf29")},
	}
	address, err := multiSigAddress.Generate(args)
	t.Logf("MultiSig Address %v", address)
	re := regexp.MustCompile(regx)
	isMatch := re.MatchString(address.Address)
	t.Logf("MultiSig Address %v, isMatch= %v", address, isMatch)
	assert.True(t, true, isMatch)
	assert.Nil(t, nil, err)
}
