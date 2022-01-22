package main

import (
	"fmt"

	"github.com/ShadiestGoat/DiscordChatExporter/config"
	"github.com/ShadiestGoat/DiscordChatExporter/fetcher"
)

const VERSION = "0.1.0"

func main() {
	fmt.Printf("Loading up Discord Channel Exporter v%v\n", VERSION)
	conf := config.Load()
	fetchConf := fetcher.ConfigType(conf)
	fetchConf.FetchMain()
}
