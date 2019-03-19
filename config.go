package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdAddr string `json:"etcd_addr"`
	MyAddr   string `json:"my_addr"`
}

var Cfg *Config


const configFile = "config.json"

func init() {
	Cfg = &Config{}
	fmt.Printf("file path %s",configFile)
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("%s:%s\n", configFile, err)
	}

	err = json.Unmarshal(buf, Cfg)

	if err != nil {
		fmt.Printf("error parsing %s:%s\n", configFile, err)
	}
}




