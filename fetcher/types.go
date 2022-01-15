package fetcher

import "github.com/ShadiestGoat/DiscordChatExporter/discord"

type JSONMetaData struct {
	Msgs []discord.Message `json:"messages"`
	IDToIndex map[string]int `json:"idToIndex"`
	ByAuthor map[string][]string `json:"byAuthor"`
	Attachments []JSONMetaAttachment `json:"attachments"`
	AuthorAttachment map[string][]int `json:"attachment_byAuthor"`
}

type JSONMetaAttachment struct {
	discord.Attachment
	AuthorID string
}

// type SingleMessageSearch struct {
// 	Msgs []discord.Message `json:"messages"`
// }