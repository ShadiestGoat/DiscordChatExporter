package fetcher

import "github.com/ShadiestGoat/DiscordChatExporter/discord"

type JSONMetaData struct {
	MsgList             []discord.Message    `json:"messages"`
	MsgIDToIndex        map[string]int       `json:"msgIdToIndex"`
	MsgByAuthor         map[string][]string  `json:"msgByAuthor"`
	AttachList     		[]JSONMetaAttachment `json:"attachments"`
	AttachIDToIndex		map[string]int		 `json:"attachmentIdToIndex"`
	AttachByAuthor		map[string][]string	 `json:"attachmentByAuthor"`
}

type JSONMetaAttachment struct {
	discord.Attachment
	AuthorID string
}
