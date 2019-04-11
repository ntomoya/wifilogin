package main

import (
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

const (
	airportPath   = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/A/Resources/airport"
	airportOption = "-I"

	premiumWi2Url = "https://service.wi2.ne.jp/wi2net/Login/1/?Wi2=1"
	tokyoTechUrl  = "https://wlanauth.noc.titech.ac.jp/login.html"
)

func airportInfo() string {
	out, err := exec.Command(airportPath, airportOption).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func currentSsid() string {
	info := airportInfo()
	lines := strings.Split(info, "\n")
	// split output by new line
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		pair := strings.Split(trim, ": ")
		if pair[0] == "SSID" {
			return pair[1]
		}
	}
	return ""
}

func loginPremiumWi2(id string, password string) int {
	timeout := 10 * time.Second
	client := http.Client{Timeout: timeout}

	data := url.Values{}
	data.Add("id", id)
	data.Add("pass", password)

	resp, err := client.PostForm(premiumWi2Url, data)
	if err != nil {
		log.Fatal(err)
	}
	return resp.StatusCode
}

func loginTokyoTech(username string, password string) int {
	timeout := 30 * time.Second
	client := http.Client{Timeout: timeout}

	data := url.Values{}
	data.Add("username", username)
	data.Add("password", password)
	data.Add("buttonClicked", "4")

	resp, err := client.PostForm(tokyoTechUrl, data)
	if err != nil {
		log.Fatal("request failed")
	}
	return resp.StatusCode
}

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	config := readConfig(filepath.Join(dir, ".config/wifilogin/config.json"))

	ssid := currentSsid()
	status := 0
	switch ssid {
	case "Wi2_club":
		status = loginPremiumWi2(config.Econnect.Id, config.Econnect.Password)
	case "TokyoTech":
		status = loginTokyoTech(config.TokyoTech.Username, config.TokyoTech.Password)
	}

	if status == 200 {
		log.Printf("login to SSID \"%s\" is successful with status code \"%d\"", ssid, status)
		return
	}
	if status != 0 {
		log.Fatalf("login to SSID \"%s\" failed with status code \"%d\"", ssid, status)
		return
	}
}
