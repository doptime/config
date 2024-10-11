package api

import "github.com/doptime/config"

// the http rpc server
type ApiSourceHttp struct {
	Name    string
	UrlBase string
	//also known as JWT token
	ApiKey string `json:"pswd"`
}

var _defaultHttpRPC = &ApiSourceHttp{Name: "doptime", UrlBase: "https://api.doptime.com", ApiKey: ""}

var APISource = []*ApiSourceHttp{_defaultHttpRPC}

func init() {
	config.LoadToml("APISource", &APISource)
}
