package main

import (
	"fmt"

	"github.com/ShadiestGoat/DiscordChatExporter/config"
	"github.com/ShadiestGoat/DiscordChatExporter/fetcher"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

const VERSION = "0.3.1"

func main() {
	fmt.Printf("Loading up Discord Channel Exporter v%v\n", VERSION)
	conf := config.Load()
	fetchConf := fetcher.ConfigType(conf)
	fetchConf.FetchMain()
	tools.Success("The downloader has finished!")
}
