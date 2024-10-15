package utils

import (
	"encoding/json"
	"fmt"
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

func LoadFromFile(key string, configObj interface{}) (err error) {
	var (
		tomlFile string = getConfigFilePath("/config.toml") // Assuming getConfigFilePath is defined elsewhere
		data            = make(map[string]interface{})
		bytes    []byte
	)

	// Read from the file and decode into the map
	if _, err = toml.DecodeFile(tomlFile, &data); err != nil {
		return fmt.Errorf("error decoding toml file: %w", err)
	}

	// Check if the key exists in the data

	if d, ok := data[key]; !ok {
		return fmt.Errorf("key %s not found in config file", key)
	} else if bytes, err = json.Marshal(d); err != nil {
		return fmt.Errorf("error encoding data into bytes: %w", err)
	} else if err = json.Unmarshal(bytes, configObj); err != nil {
		return fmt.Errorf("error decoding data into configObj: %w", err)
	}
	return nil
}

func SaveTomlFile(keyname string, configObj interface{}) (err error) {
	var (
		writer *os.File
	)
	currentConfig := map[string]interface{}{}
	//save to the file
	localConfigFile := getConfigFilePath("/config.toml")
	//read In the current configuration, and decode to currentConfig
	toml.DecodeFile(localConfigFile, &currentConfig)
	currentConfig[keyname] = configObj

	if writer, err = os.OpenFile(localConfigFile, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return
	}
	defer writer.Close()

	//write the configuration to the file
	if err = toml.NewEncoder(writer).Encode(currentConfig); err != nil {
		logger.Error().Err(err).Msg("LoadConfig_FromWeb unable to save to toml file")
	}
	return err
}
