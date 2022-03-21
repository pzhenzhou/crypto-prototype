package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	common "github.com/pzhenzhou/crypto-prototype/pkg"
	"github.com/pzhenzhou/crypto-prototype/pkg/crypto"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	errorMessageFormat = "Request parameter is invalid Name:%s, value: %s"
)

type webHandler = func(c *gin.Context)

type Response struct {
	Code    int
	Message string
	Data    interface{}
}

func responseNoData(code int, message string) Response {
	return Response{Code: code, Message: message, Data: nil}
}

func responseWithData(err error, address crypto.Address) (int, Response) {
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
	addressGeneratorCaller = crypto.AddGeneratorCaller()
	httpRouter             = map[string][]string{
		"GET": {"/check_health", "/segwit_address/:seed/:path/:password", "/multisig_address/:m/:n/:pks"},
	}

	handlerFunc = map[string]webHandler{
		"/check_health":                         checkHealth(),
		"/segwit_address/:seed/:path/:password": segWitAddressHandler(),
		"/multisig_address/:m/:n/:pks":          multiSigHandler(),
	}
	logger = common.GetLogger()
)

func HttpHandlerInit(port int) {
	if common.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	for httpMethod, pathSlices := range httpRouter {
		for _, path := range pathSlices {
			router.Handle(httpMethod, path, handlerFunc[path])
		}
	}
	var startErr = router.Run(":" + strconv.Itoa(port))
	if startErr != nil {
		logger.Error("crypto httpServer start failure", zap.Error(startErr))
		os.Exit(-1)
	}
}

func multiSigHandler() webHandler {
	return func(c *gin.Context) {
		m := c.Param("m")
		n := c.Param("n")
		pks := c.Param("pks")
		var multiPair = crypto.MultiSigNumPair{}
		if mInt, err := strconv.Atoi(m); err != nil {
			c.JSONP(http.StatusBadRequest,
				responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "MultiSig M", m)))
		} else {
			multiPair.M = mInt
		}
		if nInt, err := strconv.Atoi(n); err != nil {
			c.JSONP(http.StatusBadRequest,
				responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "MultiSig M", n)))
		} else {
			multiPair.N = nInt
		}
		if len(pks) == 0 || pks == "" {
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
		code, rsp := responseWithData(err, *address)
		c.JSONP(code, rsp)
	}
}

func segWitAddressHandler() webHandler {
	return func(c *gin.Context) {
		path := c.Param("path")
		if len(path) == 0 || path == "" {
			c.JSONP(http.StatusBadRequest, responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "path", path)))
		} else {
			if !common.IsInvalidPath(path) {
				c.JSONP(http.StatusBadRequest, responseNoData(http.StatusBadRequest, fmt.Sprintf(errorMessageFormat, "path", path)))
			}
		}

		args := map[crypto.GenerateArgs]interface{}{
			crypto.InputSeed:     c.Param("seed"),
			crypto.InputPath:     path,
			crypto.InputPassword: c.Param("password"),
		}
		address, err := addressGeneratorCaller[crypto.HDSegWitAddressGenerator].Generate(args)

		code, rsp := responseWithData(err, *address)
		c.JSONP(code, rsp)
	}
}

func checkHealth() webHandler {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "I'm Ok")
	}
}
