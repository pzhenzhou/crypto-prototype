# Crypto-prototype

### Description

This project is a simple bitcoin wallet HTTP service that implements the following features.

1. generate mnemonic and seeds according to different languages, currently, only Chinese and English are supported.
2. generate HD SegWit Address based on seeds and paths
3. generate multi-signature addresses (n-out-of-m Multisignature P2SH)

### How to build and run

#### build from source

```shell
# golang 1.18 
git clone git@github.com:pzhenzhou/crypto-prototype.git
cd crypto-prototype
# for macos intel chip
make build-darwin
# for macos apple chip
make build-darwin-arm64
# for linux
make build-linux64
```

#### run

```shell
# If you want to run in debug mode
export CRYPTO_RUN_ENV=dev
# If you want to run in prod mode
export CRYPTO_RUN_ENV=prod
# If you do not specify any arguments, the default port of the web service is 4567 and the default path to the configuration file is ./config. The system loads a list of mnemonics from the config directory
./bin/crypto-http-arm64 
# Assign the web service port and configuration file path via the command line
./bin/crypto-http-arm64 --port XXXX --config ./config 
```

### REST API

> TODO 

### Third-party lib

1. https://github.com/btcsuite/btcd
2. https://github.com/tyler-smith/go-bip32
3. https://github.com/gin-gonic/gin
