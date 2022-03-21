package main

import (
	common "github.com/pzhenzhou/crypto-prototype/pkg"
	"github.com/pzhenzhou/crypto-prototype/pkg/web"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	PortArg   string = "port"
	ConfigArg string = "config"
)

var (
	port   int
	config string
	logger = common.GetLogger()
)

func main() {
	pflag.Int(PortArg, 4567, "http server port. If not set the default is 4567")
	pflag.String(ConfigArg, "./config", "config absolute path. by default ./config")
	pflag.Parse()
	var flagErr = viper.BindPFlags(pflag.CommandLine)
	if flagErr != nil {
		logger.Error("crypto load command line arguments error", zap.Error(flagErr))
		panic(flagErr)
	}
	port = viper.GetInt(PortArg)
	config = viper.GetString(ConfigArg)
	logger.Info("crypto http service will be start", zap.Any("port", port),
		zap.Any("configPath", config))
	common.LoadWordsList(config)
	web.HttpHandlerInit(port)
}
