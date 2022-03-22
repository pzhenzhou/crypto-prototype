package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	common "github.com/pzhenzhou/crypto-prototype/pkg"
	"github.com/pzhenzhou/crypto-prototype/pkg/crypto"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

const (
	errorMessageFormat = "Request parameter is invalid Name:%s, value: %s"
)

type webHandler = func(c *gin.Context)

// Response The response object of the web service
type Response struct {
	// Similar in meaning to the return value of http status code.
	// The difference is that the http protocol represents the transport level,
	// while the code represents more of a business meaning.
	// 200: Success, 400: bad request, 500: service inner error.
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func responseNoData(code int, message string) Response {
	return Response{Code: code, Message: message, Data: nil}
}

func responseWithData(err error, address *crypto.Address) (int, Response) {
	if err != nil {
		return http.StatusInternalServerError, Response{
			http.StatusInternalServerError,
			err.Error(),
			nil,
		}
	}
	return http.StatusOK, Response{
		http.StatusOK,
		"",
		address,
	}
}

var (
	addressGeneratorCaller map[string]crypto.AddressGenerator
	httpRouter             = map[string][]string{
		"GET": {"/check_health", "/segwit_address", "/segwit_address_from_seed", "/multisig_address/:m/:n/:pks"},
	}

	handlerFunc = map[string]webHandler{
		"/check_health":                checkHealth(),
		"/segwit_address":              segWitAddressHandler(),
		"/segwit_address_from_seed":    sedWitAddressFromSeedHandler(),
		"/multisig_address/:m/:n/:pks": multiSigHandler(),
	}
	logger = common.GetLogger()
)

func HttpHandlerInit(port int) {
	if common.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	addressGeneratorCaller = crypto.AddGeneratorCaller()
	router := gin.Default()
	for httpMethod, pathSlices := range httpRouter {
		for _, path := range pathSlices {
			router.Handle(httpMethod, path, handlerFunc[path])
		}
	}
	var startErr = router.Run(":" + strconv.Itoa(port))
	if startErr != nil {
		logger.Error("crypto httpServer start failure", zap.Error(startErr))
		panic(startErr)
	}
}

func multiSigHandler() webHandler {
	return func(c *gin.Context) {
		m := c.Param("m")
		n := c.Param("n")
		pks := c.Param("pks")
		var multiPair = crypto.MultiSigNumPair{}
		if mInt, err := strconv.Atoi(m); err != nil {
			logger.Warn("MultiSig invalid request parameter", zap.Any("M", m))
			c.JSONP(http.StatusBadRequest,
				responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "MultiSig M", m)))
		} else {
			multiPair.M = mInt
		}
		if nInt, err := strconv.Atoi(n); err != nil {
			logger.Warn("MultiSig invalid request parameter", zap.Any("n", n))
			c.JSONP(http.StatusBadRequest,
				responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "MultiSig M", n)))
		} else {
			multiPair.N = nInt
		}
		if len(pks) == 0 || pks == "" {
			logger.Warn("MultiSig invalid request parameter", zap.Any("PublicKeys", pks))
			c.JSONP(http.StatusBadRequest,
				responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "MultiSig M", n)))
		}
		pksSlice := strings.Split(pks, ",")
		pksBytes := make([][]byte, len(pksSlice))
		for i, publicKey := range pksSlice {
			pksBytes[i] = []byte(publicKey)
		}
		args := map[crypto.GenerateArgs]interface{}{
			crypto.MultiSigNum:       multiPair,
			crypto.MultiSigPublicKey: pksBytes,
		}
		address, err := addressGeneratorCaller[crypto.NofMMultiSigAddressGenerator].Generate(args)
		code, rsp := responseWithData(err, address)
		c.JSONP(code, rsp)
	}
}

func sedWitAddressFromSeedHandler() webHandler {
	return func(c *gin.Context) {
		checkPath(c)
		path := strings.ReplaceAll(c.Query("path"), "\"", "")
		args := map[crypto.GenerateArgs]interface{}{
			crypto.InputSeed: c.Query("seed"),
			crypto.InputPath: path,
		}
		address, err := addressGeneratorCaller[crypto.HDSegWitAddressGenerator].Generate(args)
		code, rsp := responseWithData(err, address)
		c.JSONP(code, rsp)
	}
}

func checkPath(c *gin.Context) {
	path := c.Query("path")
	if len(path) == 0 || path == "" {
		logger.Warn("MultiSig invalid request parameter", zap.Any("path", path))
		c.JSONP(http.StatusBadRequest, responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "path", path)))
	} else {
		path = strings.ReplaceAll(path, "\"", "")
		if !common.IsInvalidPath(path) {
			logger.Warn("MultiSig invalid request parameter", zap.Any("path", path))
			c.JSONP(http.StatusBadRequest, responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "path", path)))
		}
	}
}

func segWitAddressHandler() webHandler {
	return func(c *gin.Context) {
		checkPath(c)
		args := make(map[crypto.GenerateArgs]interface{})
		args[crypto.InputPath] = strings.ReplaceAll(c.Query("path"), "\"", "")
		if len(c.Query("mnemonic")) > 0 || c.Query("mnemonic") != "" {
			args[crypto.InputMnemonic] = c.Query("mnemonic")
		}
		args[crypto.InputPassword] = c.Query("password")
		address, err := addressGeneratorCaller[crypto.HDSegWitAddressGenerator].Generate(args)
		code, rsp := responseWithData(err, address)
		c.JSONP(code, rsp)
	}
}

func checkHealth() webHandler {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "I'm Ok")
	}
}
