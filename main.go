package main

import (
	"github.com/ShadiestGoat/DiscordChatExporter/config"
	"github.com/ShadiestGoat/DiscordChatExporter/fetcher"
)

func main() {
	conf := config.Load()
	fetchConf := fetcher.ConfigType(conf)
	fetchConf.FetchAll()
}