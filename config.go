package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	EtcdAddr string `json:"etcd_addr"`
	MyAddr   string `json:"my_addr"`
}

var Cfg *Config


var configFile string

func init() {
	Cfg = &Config{}
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	configFile := path
	configFile += "/config.json"

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




