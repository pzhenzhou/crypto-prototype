## Web Service Doc

All web API return a response object in JSON format, the meaning of each field is as follows

````gotemplate
type Response struct {
	// Similar in meaning to the return value of http status code.
	// The difference is that the http protocol represents the transport level,
	// while the code represents more of a business meaning.
	// 200: Success, 400: bad request, 500: service inner error.
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
````
Recommend to use [http command line tool](https://httpie.io/) for API testing. On mac os, you can install it with the following command. The use cases in the following documentation are tested with HTTPie
```shell
brew install httpie
```

### Web API Descriptions and Samples

| HTTP Method | GET                                                          |
| ----------- | ------------------------------------------------------------ |
| URL         | /segwit_address                                              |
| REQUEST     | Query String Parameter <br> **Require**  path<br> **Option**    mnemonic , password |
| COMMENT     | If the query string in the URL does not contain a mnemonic, the system will generate a 12-digit English mnemonic |
#### Example
```shell
http get http://localhost:3456/segwit_address?mnemonic="legal winner thank year wave sausage worth useful legal winner thank yellow"&password=TREZOR&path="m/44'/0'/0'/0/0"
```
```json
{
    "code": 200,
    "data": {
        "address": "bc1q2z30tm2v0sezzc0dkrhmwqxcjeylpqxnyra3f8",
        "mnemonic": "legal winner thank year wave sausage worth useful legal winner thank yellow",
        "privateKey": "xprv9s21ZrQH143K2gA81bYFHqU68xz1cX2APaSq5tt6MFSLeXnCKV1RVUJt9FWNTbrrryem4ZckN8k4Ls1H6nwdvDTvnV7zEXs2HgPezuVccsq",
        "publicKey": "xpub6H4aDLfYSjx65SgPuBr3vcHhMPdS35JA3HJfUMLxZUHwtutUP8rHci29MEk791ZxYuxPcCFrQ8ZCRqbscGEdRMy3PzLPubNG2htkP4Niih3",
        "seed": "2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607"
    }
}
```



| HTTP Method | GET                                                          |
| ----------- | ------------------------------------------------------------ |
| URL         | /segwit_address_from_seed                                    |
| REQUEST     | Query String Parameter <br/> **Require**  seed<br/> **Require**  path |
| COMMENT     | Seed is encoded by calling method  **hex.EncodeToString(seed_byte)** to get |

#### Example
```shell
http get http://localhost:3456/segwit_address_from_seed?seed="2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607"&path="m/44'/0'/0'/0/0"
```
```json
{
    "code": 200,
    "data": {
        "address": "bc1q2z30tm2v0sezzc0dkrhmwqxcjeylpqxnyra3f8",
        "privateKey": "xprv9s21ZrQH143K2gA81bYFHqU68xz1cX2APaSq5tt6MFSLeXnCKV1RVUJt9FWNTbrrryem4ZckN8k4Ls1H6nwdvDTvnV7zEXs2HgPezuVccsq",
        "publicKey": "xpub6H4aDLfYSjx65SgPuBr3vcHhMPdS35JA3HJfUMLxZUHwtutUP8rHci29MEk791ZxYuxPcCFrQ8ZCRqbscGEdRMy3PzLPubNG2htkP4Niih3",
        "seed": "2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607"
    }
}
```



| HTTP Method | GET                                                          |
| ----------- | ------------------------------------------------------------ |
| URL         | /segwit_address_from_seed/:m/:n/:pks                         |
| REQUEST     | **Require**  m int<br>**Require**  n int <br>**Require** pks string |
| COMMENT     | Multiple pks are separated by commas. the pk in the standard Bitcoin base58 encoding |

#### Example
````shell
http get http://localhost:3456/multisig_address/3/2/020f8796e0f870a9a3b269be3b1e78e380c9b569885f0de98a9ff061c4a66e79d2,02dfa8990f3f015ff20e9b31b85ea36d47470220615fb2ac1597e20fc830727b25,03fbfbdc5df9c60e4b747805552686199e85299a5e87804dbb66a14597ddabcf29
````
```json
{
    "code": 200,
    "data": {
        "address": "bc1q2z30tm2v0sezzc0dkrhmwqxcjeylpqxnyra3f8",
        "privateKey": "xprv9s21ZrQH143K2gA81bYFHqU68xz1cX2APaSq5tt6MFSLeXnCKV1RVUJt9FWNTbrrryem4ZckN8k4Ls1H6nwdvDTvnV7zEXs2HgPezuVccsq",
        "publicKey": "xpub6H4aDLfYSjx65SgPuBr3vcHhMPdS35JA3HJfUMLxZUHwtutUP8rHci29MEk791ZxYuxPcCFrQ8ZCRqbscGEdRMy3PzLPubNG2htkP4Niih3",
        "seed": "2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607"
    }
}
```