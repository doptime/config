package settings

import (
	"github.com/doptime/config"

	"github.com/doptime/doptime/dlog"
	"github.com/rs/zerolog"
)

type ConfigSettings struct {
	//{	-1:zerolog.TraceLevel, 0: zerolog.DebugLevel, 1: zerolog.InfoLevel, 2: zerolog.WarnLevel, 3: zerolog.ErrorLevel, 4: zerolog.FatalLevel, 5: zerolog.PanicLevel, 6: zerolog.NoLevel, 7: erolog.Disabled,}
	LogLevel int8
	//super user token, this is used to bypass the security check in data access
	//SUToken is designed to allow debugging in production environment without  change the permission table permanently
	SUToken string `json:"pswd"`
}

var Settings = ConfigSettings{LogLevel: 1}

func init() {
	config.LoadToml("Settings", &Settings)
	//apply log level
	if Settings.LogLevel >= -1 && Settings.LogLevel <= 7 {
		dlog.Logger.Level(zerolog.Level(Settings.LogLevel))
	}
}
