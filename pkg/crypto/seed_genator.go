package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/pkg/errors"
	common "github.com/pzhenzhou/crypto-prototype/pkg"
	"go.uber.org/zap"
	"golang.org/x/crypto/pbkdf2"
	"strconv"
	"strings"
	"sync"
)

type SeedLen int
type WordCount int

const (
	Bit128Len SeedLen = 128
	Bit256Len SeedLen = 256

	Word12 WordCount = 12
	Word24 WordCount = 24

	passwordSalt string = "mnemonic"
)

var (
	supportLanguage        = common.SupportLanguageSlice()
	unSupportLanguageError = errors.Errorf("current only support language %v", supportLanguage)
	unSupportWordLenError  = errors.New("current only 12 or 24  mnemonic phrase")
	once                   sync.Once
	seedGeneratorInstance  *SeedGenerator
	SeedSplitError         = errors.New("Seed Split Error. binary % 11 != 0")
	mnemonicLen            = map[WordCount]SeedLen{
		Word12: Bit128Len,
		Word24: Bit256Len,
	}
	logger = common.GetLogger()
)

type SeedGenerator struct {
	bip39Word map[common.Language][]string
}

func GetSeedGenerator(words map[common.Language][]string) *SeedGenerator {
	if seedGeneratorInstance == nil {
		once.Do(func() {
			seedGeneratorInstance = &SeedGenerator{
				bip39Word: words,
			}
		})
	}
	return seedGeneratorInstance
}

// NewMnemonic
// 1. Generate a 128-bit random number and add 4 bits of checksum to the random number to get a 132-bit number
// 2. in every 11 bits to do the cut, get 12 binary numbers
// 3. Use the number generated in the 2nd step to look up the word list defined by BIP39, so as to get 12  mnemonics
// https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
func (g *SeedGenerator) NewMnemonic(input common.Language, count WordCount) (string, error) {
	if !common.IsSupportLanguage(input) {
		return "", unSupportLanguageError
	}
	if _, ok := mnemonicLen[count]; !ok {
		return "", unSupportWordLenError
	}
	if entropy, err := entropy(mnemonicLen[count]); err == nil {
		if mnemonicArray, newMnemonicErr := mnemonic(entropy, g.bip39Word[input]); newMnemonicErr != nil {
			return "", newMnemonicErr
		} else {
			return strings.Join(mnemonicArray, " "), nil
		}
	} else {
		return "", err
	}
}

func (g *SeedGenerator) NewSeed(mnemonic string, password string) []byte {
	return pbkdf2.Key([]byte(mnemonic), []byte(passwordSalt+password), 2048, 64, sha512.New)
}

func entropy(seedLen SeedLen) ([]int, error) {
	binaryString, randErr := randEntropy(seedLen)
	if randErr != nil {
		logger.Error("randEntropy() error", zap.Any("seedLen", seedLen), zap.Error(randErr))
		return nil, randErr
	}
	intSlice, err := bytesToInts(binaryString)
	if err != nil {
		logger.Error("Entropy() bytesToInts Error", zap.Any("seedLen", seedLen), zap.Error(err))
		return nil, err
	}
	return intSlice, nil
}

func mnemonic(randSlices []int, words []string) ([]string, error) {
	mnemonicSlice := make([]string, 0)
	for _, index := range randSlices {
		word := words[index]
		mnemonicSlice = append(mnemonicSlice, word)
	}
	return mnemonicSlice, nil
}

func bytesToInts(byteString string) ([]int, error) {
	byteStringLen := len(byteString)
	if byteStringLen%11 != 0 {
		return nil, SeedSplitError
	}
	var randSlice = make([]int, 0)
	len := byteStringLen / 11
	for i := 0; i < len; i++ {
		start := 11 * i
		end := (11 * i) + 11
		strVal := byteString[start:end]
		intVal, err := strconv.ParseInt(strVal, 2, 32)
		if err != nil {
			logger.Error("string convert int32 error", zap.Any("stringValue", strVal), zap.Error(err))
			return nil, err
		}
		randSlice = append(randSlice, int(intVal))
	}
	return randSlice, nil
}

func checkSumBinary(binarySlices []byte, seedLen SeedLen) (string, error) {
	sha256Bytes := sha256.Sum256(binarySlices)
	len := int(seedLen) / 32
	var buff bytes.Buffer
	for _, byteValue := range sha256Bytes {
		binary := fmt.Sprintf("%08b", byteValue)
		if len <= 8 {
			buffValue := binary[:len]
			buff.WriteString(buffValue)
			break
		} else {
			buff.WriteString(binary)
			len -= 8
		}
	}
	return buff.String(), nil
}

func bytesEncode(byteSlice []byte) string {
	var strBuffer bytes.Buffer
	for _, byteElement := range byteSlice {
		binaryStr := fmt.Sprintf("%08b", byteElement)
		strBuffer.WriteString(binaryStr)
	}
	return strBuffer.String()
}

func randEntropy(seedLen SeedLen) (string, error) {
	entropyByteLen := int(seedLen / 8)
	byteSlice := make([]byte, entropyByteLen)
	_, err := rand.Read(byteSlice)
	if err != nil {
		logger.Error("generate crypto.rand.Read() error cause by", zap.Error(err))
		return "", err
	}

	encodeValue := bytesEncode(byteSlice)
	checkSumValue, err := checkSumBinary(byteSlice, seedLen)
	if err != nil {
		logger.Error("checkSumBinary error cause by", zap.Error(err))
		return "", err
	}
	bitsRs := encodeValue + checkSumValue
	return bitsRs, nil
}
