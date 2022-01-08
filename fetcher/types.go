package fetcher

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

type Message struct {
	ID string `json:"id"`
	Content string `json:"content"`
	Channel string `json:"channel_id"`
	Author Author `json:"author"`
	Attachments []Attachment `json:"attachments"`
	Embeds []Embed `json:"embeds"`
	Mentions []Author `json:"mentions"`
	MentionRoles []string `json:"mention_roles"`
	MentionsEveryone bool `json:"mention_everyone"`
	Timestamp int `json:"timestamp"`
	IsEdited bool
}

type Author struct {
	ID string `json:"id"`
	Name string `json:"username"`
	Avatar string `json:"avatar"`
}

type Attachment struct {
	ID string `json:"id"`
	Name string `json:"filename"`
	Size int `json:"size"`
	Url string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
	ContentType string `json:"content_type"`
}

type EmbedType int8

const (
	EMBED_LINK = iota
	EMBED_GIF
	EMBED_IGNORE
)

type Embed struct {
	Type EmbedType
	Url string
	Title string
	Description string
	Color string
}

type EmbedVideo struct {
	Url string `json:"url"`
}

func (msg *Message) UnmarshalJSON(b []byte) error {
	var s map[string]json.RawMessage
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	for key, value := range s {
		switch key {
			case "id":
				parsed := ""
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.ID = parsed

				i, err := strconv.ParseInt(parsed, 10, 64)
				tools.PanicIfErr(err)

				timestamp := (i >> 22) + 1420070400000
				
				msg.Timestamp = int(time.Unix(0, timestamp*1000000).UnixMicro())
			case "content":
				parsed := ""
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.Content = parsed
			case "channel_id":
				parsed := ""
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)
				msg.Channel = parsed
			case "author":
				parsed := Author{}
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.Author = parsed
			case "attachments":
				parsed := []Attachment{}
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.Attachments = parsed
			case "embeds":
				parsed := []Embed{}
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.Embeds = parsed
			case "mentions":
				parsed := []Author{}
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.Mentions = parsed
			case "mention_roles":
				parsed := []string{}
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.MentionRoles = parsed
			case "mention_everyone":
				parsed := false
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)
			
				msg.MentionsEveryone = parsed
			case "edited_timestamp":
				msg.IsEdited = string(value) != "null"
		}
	}
	return nil
}

func (embed *Embed) UnmarshalJSON(b []byte) error {
	var s map[string]json.RawMessage
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}


	typePre := ""

	err := json.Unmarshal(s["type"], &typePre)
	
	tools.PanicIfErr(err)

	switch typePre {
		case "gifv":
			embed.Type = EMBED_GIF
		case "link", "rich":
			embed.Type = EMBED_LINK
		default:
			embed.Type = EMBED_IGNORE
	}

	for key, value := range s {
		switch key {
			case "type":
			case "url":
				if embed.Type == EMBED_GIF {
					continue
				}
				parsed := ""
				
				err = json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				embed.Url = parsed
			case "title":
				parsed := ""
				
				err = json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				embed.Title = parsed
			case "description":
				parsed := ""
				err = json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)
				embed.Description = parsed
			case "color":
				parsedInt, err := strconv.ParseInt(string(value), 10, 64)
				tools.PanicIfErr(err)
				embed.Color = strconv.FormatInt(parsedInt, 16)
			case "video":
				if embed.Type != EMBED_GIF {
					continue
				}
				parsed := EmbedVideo{}
				
				err = json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				embed.Url = parsed.Url
		}
	}
	return nil
}


type JSONMetaData struct {
	Msgs []Message `json:"messages"`
	IDToIndex map[string]int `json:"idToIndex"`
	ByAuthor map[string][]string `json:"byAuthor"`
	Attachments []JSONMetaAttachment `json:"attachments"`
	AuthorAttachment map[string]int `json:"attachment_byAuthor"`
}

type JSONMetaAttachment struct {
	Attachment
	AuthorID string
}

// 	"mentions": [],
// 	"mention_roles": [],
// 	"pinned": false,
// 	"mention_everyone": false,
// 	"tts": false,
// 	"timestamp": "2022-01-06T21:56:07.979000+00:00",
// 	"edited_timestamp": null,
// 	"flags": 0,
// 	"components": [],
// 	"hit": true
