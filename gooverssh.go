package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/scottkiss/gosshtool"
	"io/ioutil"
	"log"
)

var config tomlConfig

func loadconfig() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		panic("load config error")
	}
}

func main() {
	loadconfig()
	server := new(gosshtool.LocalForwardServer)
	server.LocalBindAddress = config.SSH.LocalBindAddress
	server.RemoteAddress = config.Remote.Host + ":" + config.Remote.Port
	server.SshServerAddress = config.SSH.SshServerAddress
	if config.SSH.SshUserPassword == "" && config.SSH.SshPrivateKey == "" {
		panic("password or private key path is empty")
	}
	if config.SSH.SshUserPassword != "" && config.SSH.SshPrivateKey == "" {
		server.SshUserPassword = config.SSH.SshUserPassword
	}
	if config.SSH.SshUserPassword == "" && config.SSH.SshPrivateKey != "" {
		buf, _ := ioutil.ReadFile(config.SSH.SshPrivateKey)
		server.SshPrivateKey = string(buf)
	}
	server.SshUserName = config.SSH.SshUserName
	server.Start(started)
	defer server.Stop()
}

func started() {
	log.Println("go over ssh server started")
}

type tomlConfig struct {
	Title  string
	SSH    ssh
	Remote remote
}

type ssh struct {
	LocalBindAddress string `toml:"local_bind_address"`
	SshServerAddress string `toml:"ssh_server"`
	SshUserName      string `toml:"ssh_user"`
	SshUserPassword  string `toml:"ssh_password"`
	SshPrivateKey    string `toml:"private_key_path"`
}

type remote struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}
