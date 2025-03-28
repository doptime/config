package config

import (
	"github.com/doptime/config/cfgurl"
	"github.com/doptime/config/utils"
	"github.com/doptime/logger"
)

func LoadItemFromToml(keyname string, configObj interface{}) {
	logger.Info().Str("Config Item", keyname).Str("data", "loading ..").Send()
	//step1: load config from file
	utils.LoadFromFile(keyname, configObj)
	//step2: load config from env. this will overwrite the config from file
	//step3: load config from web. this will overwrite the config from env.
	//warning local config will be overwritten by the config from web, to prevent falldown of config from web.
	cfgurl.LoadFromUrl(keyname, configObj)
	logger.Info().Str("Config Item", keyname).Str("data", utils.ToHidePswdString(configObj)).Send()
}
