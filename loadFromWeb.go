package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/doptime/logger"
)

var CfgUrl = struct {
	ConfigUrl string `toml:"CONFIG_URL"`
}{}

func getConfigUrl() (configUrl string) {
	if CfgUrl.ConfigUrl != "" {
		return CfgUrl.ConfigUrl
	}
	//env first: try read from env
	CfgUrl.ConfigUrl = os.Getenv("CONFIG_URL")
	//secondary option: if not in env read from file
	if CfgUrl.ConfigUrl == "" && loadFromFile(&CfgUrl) != nil {
		return ""
	}
	saveTomlFile("CONFIG_URL", CfgUrl.ConfigUrl)
	return CfgUrl.ConfigUrl
}

var getCachedTomlFromUri = func() func() (page string, err error) {
	var configUrl string = os.Getenv("CONFIG_URL")
	var LastLoadTime time.Time = time.Now().Add(-time.Hour * 24 * 365)
	var LastCachedPage string = ""

	type ConfigUrl struct {
		ConfigUrl string
	}

	//read from env as primary source, and then from the configuration file as secondary source
	if configUrl == "" {
		var cfgurl ConfigUrl
		err := loadFromFile(&cfgurl)
		if err != nil {
			configUrl = CfgUrl.ConfigUrl
		}
	}
	downloadPage := func() (page string, err error) {
		//return if the url is not valid or empty
		if !strings.HasPrefix(strings.ToLower(configUrl), "http") {
			return "", fmt.Errorf("invalid CONFIG_URL")
		}
		var resp *http.Response
		var pageBytes []byte
		//download from the url and save to the file
		httpClient := &http.Client{Timeout: time.Second * 6}
		if resp, err = httpClient.Get(configUrl); err != nil {
			err = fmt.Errorf("failed to download CONFIG_URL page: " + err.Error())
			logger.Error().Err(err).Str("Url", configUrl).Msg("LoadConfig_FromWeb failed")
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			logger.Error().Str("Url", configUrl).Int("StatusCode", resp.StatusCode).Msg("LoadConfig_FromWeb failed")
			return "", fmt.Errorf("failed to download CONFIG_URL page: " + resp.Status)
		}
		//read the page
		pageBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("failed to read CONFIG_URL page: " + err.Error())
			logger.Error().Err(err).Str("Url", configUrl).Msg("LoadConfig_FromWeb failed")
			return "", err
		}
		return string(pageBytes), nil
	}

	return func() (page string, err error) {
		if configUrl == "" {
			return "", fmt.Errorf("empty CONFIG_URL")
		}
		pageExpired := time.Now().After(LastLoadTime.Add(time.Minute * 5))
		if pageExpired {
			LastCachedPage, err = downloadPage()
			if err != nil {
				return "", err
			}
			LastLoadTime = time.Now()
		}
		return LastCachedPage, nil
	}
}()

func loadFromUrl(configUrl string, keyname string, configObj interface{}) (err error) {
	var (
		page string
	)
	if page, err = getCachedTomlFromUri(); err != nil {
		return err
	} else if page == "" {
		return fmt.Errorf("empty CONFIG_URL page")
	}

	//decode the page to the configuration object
	if _, err = toml.Decode(page, configObj); err != nil {
		logger.Error().Err(err).Str("Url", configUrl).Msg("LoadConfig_FromWeb failed")
		return
	}
	return saveTomlFile(keyname, configObj)
}

func saveTomlFile(keyname string, configObj interface{}) (err error) {
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
