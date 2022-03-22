package common

import (
	"bufio"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type CryptoEnv string
type Language string

const (
	EnvPrefix         CryptoEnv = "CRYPTO"
	RunEnv            CryptoEnv = "RUN_ENV"
	English           Language  = "english"
	ChineseSimplified Language  = "chinese_simplified"
)

func (c CryptoEnv) String() string {
	return string(c)
}

var (
	EnvSlice        = []CryptoEnv{RunEnv}
	bitcoinPropose  = []int{44, 49, 84}
	logger          *zap.Logger
	supportLanguage = map[Language]bool{
		English:           true,
		ChineseSimplified: true,
	}
	wordList = map[Language][]string{}
)

func GetLogger() *zap.Logger {
	return logger
}

func SupportLanguageSlice() []Language {
	rs := make([]Language, 0)
	for language, support := range supportLanguage {
		if support {
			rs = append(rs, language)
		}
	}
	return rs
}

func IsSupportLanguage(input Language) bool {
	for language, support := range supportLanguage {
		if strings.EqualFold(string(language), string(input)) && support {
			return true
		}
	}
	return false
}

func GetWordList() map[Language][]string {
	var res = make(map[Language][]string)
	for key, value := range wordList {
		res[key] = value
	}
	return res
}

func init() {
	fmt.Println("crypto init bind ENV. env names = ", EnvSlice)
	viper.SetEnvPrefix(EnvPrefix.String())
	viper.AllowEmptyEnv(true)
	for _, envName := range EnvSlice {
		var bindErr = viper.BindEnv(envName.String())
		if bindErr != nil {
			fmt.Println("crypto viper.BindEnv error cause by ", bindErr)
			os.Exit(-1)
		}
	}
	initLogger()
}

func LoadWordsList(configPath string) {
	files, err := ioutil.ReadDir(configPath)
	if err != nil {
		fmt.Println("Read ConfigDir Error", err)
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		filenameWithoutExt := strings.Replace(fileName, path.Ext(fileName), "", -1)
		if val, ok := supportLanguage[Language(filenameWithoutExt)]; ok && val {
			filePtr, _ := os.Open(configPath + "/" + fileName)
			scanner := bufio.NewScanner(filePtr)
			scanner.Split(bufio.ScanWords)
			wordBuf := make([]string, 0)
			for scanner.Scan() {
				if scanner.Text() == "" {
					continue
				}
				wordBuf = append(wordBuf, strings.TrimSpace(scanner.Text()))
			}
			wordList[Language(filenameWithoutExt)] = wordBuf
		}
	}
}

func IsProd() bool {
	var runEnvValue = viper.Get(RunEnv.String())
	if runEnvValue == nil {
		runEnvValue = "dev"
	}
	fmt.Println("RunEnv = ", runEnvValue)
	return strings.EqualFold(runEnvValue.(string), "prod")
}

func initLogger() {
	var zapLogger *zap.Logger
	var logInitErr error
	if IsProd() {
		zapLogger, logInitErr = zap.NewProduction()
	} else {
		zapLogger, logInitErr = zap.NewDevelopment()
	}
	if logInitErr != nil {
		fmt.Println("crypto Init Logger Error cause by", logInitErr)
		os.Exit(-1)
	}
	logger = zapLogger
}

// IsInvalidPath m/purpose'/coin'/account'/ change/address_index
func IsInvalidPath(path string) bool {
	if len(strings.Split(path, "/")) > 2 {
		indexSlice := strings.Split(path, "/")[1]
		child := GetChild(indexSlice)
		for _, code := range bitcoinPropose {
			if child == uint32(code)+0x80000000 {
				return true
			}
		}
	}
	return false
}

func GetChild(index string) uint32 {
	if index[len(index)-1] == 39 {
		return cast.ToUint32(index[:len(index)-1]) + 0x80000000
	}
	return cast.ToUint32(index)
}
