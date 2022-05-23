package master

import (
	"encoding/json"
	"io/ioutil"
)

// 程序配置
type Config struct {
	ApiPort         int `json:"apiPort"`
	ApiReadTimeOut  int `json:"apiReadTimeout"`
	ApiWriteTimeOut int `json:"apiWriteTimeOut"`
}

var G_config *Config

func InitConfig(filename string) (err error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	conf := Config{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return
	}
	G_config = &conf
	return
}
