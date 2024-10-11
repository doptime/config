package http

import (
	"github.com/doptime/config"
)

type ConfigHttp struct {
	CORES string
	Port  int64  `env:"Port,default=80"`
	Path  string `env:"Path,default=/"`
	//MaxBufferSize is the max size of a task in bytes, default 10M
	MaxBufferSize int64  `env:"MaxBufferSize,default=10485760"`
	JwtSecret     string `env:"Secret"  json:"pswd"`
	//AutoAuth should never be true in production
	AutoAuth bool `env:"AutoAuth,default=false"`
}

var Http ConfigHttp = ConfigHttp{CORES: "*", Port: 80, Path: "/", MaxBufferSize: 10485760}

func init() {
	config.LoadToml("Http", &Http)
}
