package cfglog

import (
	"github.com/doptime/config"

	"github.com/doptime/doptime/dlog"
	"github.com/rs/zerolog"
)

type Log struct {
	//{	-1:zerolog.TraceLevel, 0: zerolog.DebugLevel, 1: zerolog.InfoLevel, 2: zerolog.WarnLevel, 3: zerolog.ErrorLevel, 4: zerolog.FatalLevel, 5: zerolog.PanicLevel, 6: zerolog.NoLevel, 7: erolog.Disabled,}
	LogLevel int8
}

var option = Log{LogLevel: 1}
var LogLevel int8 = 1

func init() {
	config.LoadToml("Log", &option)
	//apply log level
	if option.LogLevel >= -1 && option.LogLevel <= 7 {
		LogLevel = option.LogLevel
		dlog.Logger.Level(zerolog.Level(option.LogLevel))
	}
}
