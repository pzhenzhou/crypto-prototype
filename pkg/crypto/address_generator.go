package crypto

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	common "github.com/pzhenzhou/crypto-prototype/pkg"
	"github.com/tyler-smith/go-bip32"
	"strings"
)

type GenerateArgs string

const (
	InputPassword                GenerateArgs = "password"
	InputSeed                    GenerateArgs = "seed"
	InputPath                    GenerateArgs = "path"
	MultiSigNum                  GenerateArgs = "multiSigPair"
	MultiSigPublicKey            GenerateArgs = "multiSigPublicKeys"
	HDSegWitAddressGenerator                  = "HDSegWitAddressGenerator"
	NofMMultiSigAddressGenerator              = "NofMMultiSigAddressGenerator"
)

var (
	ArgsMustBeNotNull        = errors.New("Input Args must be not null")
	PasswordInvalid          = errors.New("Invalid password, Password must not be empty and len(password) >= 6")
	MultiSigArgsInvalid      = errors.New("n-out-of-m MultiSig argument must not be empty.")
	MultiSigNumValueInvalid  = errors.New("n-out-of-m MultiSig.invalid N or M")
	MultiSigPublicKeyInvalid = errors.New("n-out-of-m MultiSig.invalid public key")
	seedGenerator            = GetSeedGenerator(common.GetWordList())
)

type MultiSigNumPair struct {
	N int
	M int
}

type Address struct {
	Address    string `json:"address"`
	PublicKey  string `json:"publicKey,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
	Mnemonic   string `json:"Mnemonic,omitempty"`
}

type AddressGenerator interface {
	Generate(args map[GenerateArgs]interface{}) (*Address, error)
}

func AddGeneratorCaller() map[string]AddressGenerator {
	return map[string]AddressGenerator{
		HDSegWitAddressGenerator:     NewHDSegWitAddress(seedGenerator),
		NofMMultiSigAddressGenerator: MultiSigAddress{},
	}
}

func extractKeyForBIP32(children []string, parent *bip32.Key) (*bip32.Key, error) {
	child, err := parent.NewChildKey(common.GetChild(children[0]))
	if err != nil {
		return nil, err
	}

	if len(children) == 1 {
		return child, nil
	}

	return extractKeyForBIP32(children[1:], child)
}

type HDSegWitAddress struct {
	seedGenerator *SeedGenerator
}

func NewHDSegWitAddress(seedGenerator *SeedGenerator) HDSegWitAddress {
	return HDSegWitAddress{
		seedGenerator: seedGenerator,
	}
}

// Generate Produce HD SegWit address based on the given mnemonic and password
// If the mnemonic is empty, the method automatically generates a 12-digit English mnemonic
// The password length cannot be less than 6
func (h HDSegWitAddress) Generate(args map[GenerateArgs]interface{}) (*Address, error) {
	if args == nil || len(args) == 0 {
		return nil, ArgsMustBeNotNull
	}
	if pwd, ok := args[InputPassword]; !ok {
		return nil, PasswordInvalid
	} else {
		if len(pwd.(string)) < 6 {
			return nil, PasswordInvalid
		}
	}
	path := args[InputPath].(string)
	password := args[InputPassword].(string)
	var seed []byte
	mnemonic := args[InputSeed]
	if inputMnemonic, ok := args[InputSeed]; !ok {
		newMnemonic, err := h.seedGenerator.NewMnemonic(common.English, Word12)
		if err != nil {
			return nil, err
		}
		mnemonic = newMnemonic
		seed = h.seedGenerator.NewSeed(newMnemonic, password)
	} else {
		seed = h.seedGenerator.NewSeed(inputMnemonic.(string), password)
	}
	masterPrivateKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	children := strings.Split(path, "/")[1:]
	bip32Key, err := extractKeyForBIP32(children, masterPrivateKey)
	witness := btcutil.Hash160(bip32Key.PublicKey().Key)
	addressHash, err := btcutil.NewAddressWitnessPubKeyHash(witness, &chaincfg.MainNetParams)

	if err != nil {
		return nil, err
	}
	return &Address{
		addressHash.EncodeAddress(),
		bip32Key.PublicKey().B58Serialize(),
		masterPrivateKey.B58Serialize(),
		mnemonic.(string),
	}, nil
}

type MultiSigAddress struct {
}

func (m MultiSigAddress) Generate(args map[GenerateArgs]interface{}) (*Address, error) {
	var multiSig MultiSigNumPair
	if _, ok := args[MultiSigNum]; !ok {
		return nil, MultiSigArgsInvalid
	} else {
		multiSig = args[MultiSigNum].(MultiSigNumPair)
	}

	if multiSig.N > 16 || multiSig.N < 1 {
		return nil, MultiSigNumValueInvalid
	}
	if multiSig.N > multiSig.M || multiSig.M < 1 {
		return nil, MultiSigNumValueInvalid
	}
	if _, ok := args[MultiSigPublicKey]; !ok {
		return nil, MultiSigPublicKeyInvalid
	}

	publicKeys := args[MultiSigPublicKey].([][]byte)

	if len(publicKeys) != multiSig.M {
		return nil, MultiSigPublicKeyInvalid
	}
	scriptBuilder := txscript.NewScriptBuilder()
	scriptBuilder.AddOp(byte(0x50 + multiSig.N))
	// add the public keys
	for _, public := range publicKeys {
		scriptBuilder.AddData(public)
	}
	scriptBuilder.AddOp(byte(0x50 + multiSig.M))
	// add the check-multi-sig OP_CODE
	scriptBuilder.AddOp(txscript.OP_CHECKMULTISIG)
	script, err := scriptBuilder.Script()
	if err != nil {
		return nil, err
	}
	redeemHash := btcutil.Hash160(script)
	address, err := btcutil.NewAddressScriptHashFromHash(redeemHash, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return &Address{
		address.EncodeAddress(),
		"",
		"",
		"",
	}, nil
}
