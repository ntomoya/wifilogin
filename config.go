package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type econnectConfig struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type tokyoTechConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type config struct {
	Econnect  econnectConfig  `json:"econnect"`
	TokyoTech tokyoTechConfig `json:"tokyotech"`
}

func readConfig(path string) *config {
	// read json file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	config := new(config)
	if err := json.Unmarshal(bytes, config); err != nil {
		log.Fatal(err)
	}

	return config
}
