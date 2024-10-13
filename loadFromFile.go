package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/doptime/logger"
)

func getConfigFilePath(filename string) string {
	var (
		tomlFilePath string
		err          error
	)
	//tomlPath is same path as the binary
	if tomlFilePath, err = os.Executable(); err != nil {
		logger.Panic().Msg("Failed to get executable path")
	}
	tomlFilePath = filepath.Dir(tomlFilePath)
	return tomlFilePath + filename
}

func loadFromFile(configObj interface{}) (err error) {

	var (
		tomlFile string = getConfigFilePath("/config.toml")
	)

	_, err = toml.DecodeFile(tomlFile, configObj)
	return err
}
