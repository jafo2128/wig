package main

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	WsHost      string `json:host`
	WsPort      int    `json:port`
	SslCert     string `json:sslcert`
	SslKey      string `json:sslkey`
	AutoGenCert bool   `json:autogencert`
	AnonAuth    bool   `json:anonauth`
}

var _cfg_file = "options.json"

func LoadSettings() *Settings {
	ret := &Settings{}
	buf, ber := ioutil.ReadFile(_cfg_file)
	if ber == nil {
		je := json.Unmarshal(buf, ret)
		if je == nil {
			return ret
		}
	} else {
		ret = &Settings{
			WsHost:      "localhost",
			WsPort:      9002,
			AutoGenCert: true,
			SslCert:     "cert.crt",
			SslKey:      "cert.key",
			AnonAuth:    true,
		}
		js, jer := json.MarshalIndent(ret, "\t", "")
		if jer == nil {
			ioutil.WriteFile(_cfg_file, js, 0644)
		}
	}
	return ret
}
