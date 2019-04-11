package main

import (
	"fmt"
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

func airportInfo() (string, error) {
	out, err := exec.Command(airportPath, airportOption).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func currentSsid() string {
	info, err := airportInfo()
	if err != nil {
		// FIXME: Do not logging here
		log.Fatal(err)
	}
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

func loginPremiumWi2(id string, password string) (int, error) {
	timeout := 10 * time.Second
	client := http.Client{Timeout: timeout}

	data := url.Values{}
	data.Add("id", id)
	data.Add("pass", password)

	resp, err := client.PostForm(premiumWi2Url, data)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func loginTokyoTech(username string, password string) (int, error) {
	timeout := 30 * time.Second
	client := http.Client{Timeout: timeout}

	data := url.Values{}
	data.Add("username", username)
	data.Add("password", password)
	data.Add("buttonClicked", "4")

	resp, err := client.PostForm(tokyoTechUrl, data)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	config := readConfig(filepath.Join(dir, ".config/wifilogin/config.json"))

	var status int
	var err error
	ssid := currentSsid()
	switch ssid {
	case "Wi2_club":
		notify("Attempt to login", fmt.Sprintf("Attempt to login SSID \"%s\"", ssid))
		status, err = loginPremiumWi2(config.Econnect.Id, config.Econnect.Password)
	case "TokyoTech":
		notify("Attempt to login", fmt.Sprintf("Attempt to login SSID \"%s\"", ssid))
		status, err = loginTokyoTech(config.TokyoTech.Username, config.TokyoTech.Password)
	}

	if err != nil {
		notify("Login Failed", fmt.Sprintf("HTTP request for login to SSID \"%s\" failed", ssid))
		log.Fatal(err)
	}
	if status == 200 {
		notify("Login Succeeded!", fmt.Sprintf("Login to SSID \"%s\" succeeded!", ssid))
		log.Printf("login to SSID \"%s\" was successful with status code \"%d\"", ssid, status)
		return
	}
	if status != 0 {
		notify("Login Failed", fmt.Sprintf("Login to SSID \"%s\" failed.", ssid))
		log.Fatalf("login to SSID \"%s\" failed with status code \"%d\"", ssid, status)
	}
}
