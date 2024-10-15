package cfgapi

import (
	"github.com/doptime/config"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// the http rpc server
type ApiSourceHttp struct {
	Name    string
	UrlBase string
	//also known as JWT token
	ApiKey string `psw:"true"`
}

var _defaultHttpRPC = &ApiSourceHttp{Name: "doptime", UrlBase: "https://api.doptime.com", ApiKey: ""}

var apiSources = []*ApiSourceHttp{_defaultHttpRPC}
var Servers cmap.ConcurrentMap[string, *ApiSourceHttp] = cmap.New[*ApiSourceHttp]()

func init() {
	config.LoadToml("APISource", &apiSources)
	Servers.Set(_defaultHttpRPC.Name, _defaultHttpRPC)
	for _, api := range apiSources {
		Servers.Set(api.Name, api)
	}
}
