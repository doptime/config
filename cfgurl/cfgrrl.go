package cfgurl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/doptime/config/utils"
	"github.com/doptime/logger"
)

// read from env as primary source, and then from the configuration file as secondary source
var getConfigUrl = func() func() string {
	var configUrl string
	var loaded bool = false

	return func() (url string) {
		if loaded {
			return configUrl
		}
		loaded = true
		//env first: try read from env
		configUrl = os.Getenv("CONFIG_URL")
		//secondary option: if not in env read from file
		if configUrl == "" {
			utils.LoadFromFile("CONFIG_URL", &configUrl)
		}
		return configUrl
	}
}()

var getCachedTomlFromUri = func() func() (page string, err error) {
	var configUrl string = getConfigUrl()
	var LastLoadTime time.Time = time.Now().Add(-time.Hour * 24 * 365)
	var LastCachedPage string = ""

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

func LoadFromUrl(keyname string, configObj interface{}) (err error) {
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
		logger.Error().Err(err).Str("Url", getConfigUrl()).Msg("LoadConfig_FromWeb failed")
		return
	}
	return utils.SaveTomlFile(keyname, configObj)
}
