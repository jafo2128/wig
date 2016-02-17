package main

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	FcgiAddress string `json:fcgiAddr`
	WsAddress   string `json:wsAddr`
}

var _cfg_file = "options.conf"

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
			FcgiAddress: ":9001",
			WsAddress:   ":9002",
		}
		js, jer := json.Marshal(ret)
		if jer == nil {
			ioutil.WriteFile(_cfg_file, js, 0644)
		}
	}
	return ret
}
