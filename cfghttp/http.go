package cfghttp

import (
	"github.com/doptime/config"
)

type ConfigHttp struct {
	CORES string
	Port  int64
	Path  string
	//MaxBufferSize is the max size of a task in bytes, default 10M
	MaxBufferSize int64
	JwtSecret     string `psw:"true"`
	//AutoAuth should never be true in production
	DangerousAutoWhitelist bool

	// super user token, this is used to bypass the security check in data access
	// SUToken is designed to allow debugging in production environment without  change the permission table permanently
	SUToken string `psw:"true"`
}

var CORES = "*"
var Port int64 = 80
var Path = "/"
var MaxBufferSize = int64(10485760)
var JWTSecret string = ""
var DangerousAutoWhitelist bool = false
var SUToken string = ""

func init() {
	var httpOption ConfigHttp
	config.LoadItemFromToml("Http", &httpOption)
	CORES = httpOption.CORES
	Port = httpOption.Port
	Path = httpOption.Path
	MaxBufferSize = httpOption.MaxBufferSize
	JWTSecret = httpOption.JwtSecret
	DangerousAutoWhitelist = httpOption.DangerousAutoWhitelist
	SUToken = httpOption.SUToken
}
