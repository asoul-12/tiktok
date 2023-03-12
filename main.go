package main

import (
	"tiktok/config"
	"tiktok/global"
	"tiktok/serverInit"
)

func main() {
	global.Config = config.LoadConfigFromFile("config-test.json")
	serverInit.ServerInitAndStart()
}
