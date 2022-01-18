package discord

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

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

				msg.Timestamp = IDToTimestamp(parsed)
			case "content":
				parsed := ""
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.Content = parsed
			case "type":
				parsed := MsgType(0)
				
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				msg.Type = parsed
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
			case "referenced_message":
				msg.IsReply = true
				parsed := ReplyMsg{}
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)
				msg.ReplyTo = parsed
			case "sticker_items":
				msg.IsSticker = true
				parsed := []Sticker{}
				err := json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)
				msg.Stickers = parsed
		}
	}

	if msg.Type != 0 && msg.Type != 19 {
		switch msg.Type {
			case 1:
				msg.Content = fmt.Sprintf(`*%v* has added someone!`, msg.Author.Name)
			case 2:
				msg.Content = fmt.Sprintf(`*%v* has removed someone!`, msg.Author.Name)
			case 3:
				msg.Content = fmt.Sprintf(`*%v* started a call!`, msg.Author.Name)
			case 4:
				msg.Content = fmt.Sprintf(`*%v* renamed the gc`, msg.Author.Name)
			default:
				panic(fmt.Sprintf("%#v", msg))
		}

		msg.Author = Author{
			Name: "System",
			Avatar: "c6a249645d46209f337279cd2ca998c7",
			ID: "643945264868098049",
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
		case "image":
			embed.Type = EMBED_IMAGE
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
			case "thumbnail":
				if embed.Type != EMBED_IMAGE {
					continue
				}
				parsed := EmbedImageThumbnail{}
				
				err = json.Unmarshal(value, &parsed)
				tools.PanicIfErr(err)

				embed.Thumbnail = parsed
		}
	}
	return nil
}