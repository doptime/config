package config

import (
	util "config/utils"

	"github.com/doptime/doptime/dlog"
)

func LoadToml(keyname string, configObj interface{}) {
	dlog.Info().Msg("Loading configuration Item \"Http\" ..")
	//step1: load config from file
	loadFromFile(configObj)
	var configUrl = getConfigUrl()
	//step2: load config from env. this will overwrite the config from file
	//step3: load config from web. this will overwrite the config from env.
	//warning local config will be overwritten by the config from web, to prevent falldown of config from web.
	loadFromUrl(configUrl, keyname, configObj)
	dlog.Info().Str("Config load", keyname).Str("json", util.ToHidePswdString(configObj)).Send()
}
