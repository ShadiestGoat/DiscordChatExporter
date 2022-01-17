package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ShadiestGoat/DiscordChatExporter/components"
	"github.com/ShadiestGoat/DiscordChatExporter/config"
	"github.com/ShadiestGoat/DiscordChatExporter/discord"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

type ConfigType config.Config

type NewDMChannelResp struct {
	Id string `json:"id"`
}

func (conf ConfigType) checkToken() {
	resp := []byte{}
	err := conf.discordFetch("/users/@me", &resp)
	if errors.Is(err, ErrBadAuth) {
		panic("Warning! This token is invalid!")
	} else {
		tools.PanicIfErr(err)
	}
	user := discord.Author{}
	json.Unmarshal(resp, &user)
	fmt.Printf("Token is valid! Logged in as %v#%v!\n", user.Name, user.Discriminator)
}

func (conf ConfigType) FetchMain() {
	conf.checkToken()

	ParsedMessages := map[string]JSONMetaData{}
	NumLeft := conf.Filter.NumMax

	maxTime := conf.Filter.MaxTime
	minTime := conf.Filter.MinTime
	maxMsg := discord.Message{}
	minMsg := discord.Message{}

	ext := "."

	switch conf.ExportType {
	case config.EXPORT_TYPE_HTML:
		ext += "html"
	case config.EXPORT_TYPE_JSON:
		ext += "json"
	case config.EXPORT_TYPE_TEXT:
		ext += "log"
	default:
		panic("Default case achieved, unknown export type!")
	}

	switch conf.IdType {
	case config.ID_TYPE_CHANNEL:
	case config.ID_TYPE_GUILD:
		parsedGuilds := map[string]bool{}

		for _, id := range conf.Ids {
			if _, ok := parsedGuilds[id]; ok {
				continue
			}
			parsedGuilds[id] = true
		}

		conf.Ids = []string{}
		for _, guildId := range parsedGuilds {
			resp := []byte{}
			err := conf.discordFetch(fmt.Sprintf("/guilds/%v/channels", guildId), &resp)
			if errors.Is(err, Err404) {
				fmt.Printf("Warning! %v is not a guild id! Ignoring...", guildId)
				continue
			} else {
				tools.PanicIfErr(err)
			}
			allChannels := []discord.Channel{}
			json.Unmarshal(resp, &allChannels)
			chanIds := []string{}
			for _, channel := range allChannels {
				switch channel.Type {
				case discord.CHANNEL_TYPE_GUILD_CATEGORY,
					discord.CHANNEL_TYPE_GUILD_STAGE_VOICE,
					discord.CHANNEL_TYPE_GUILD_STORE,
					discord.CHANNEL_TYPE_GUILD_VOICE:
					continue
				}
				if channel.Nsfw && conf.IgnoreNsfw {
					continue
				}
				chanIds = append(chanIds, channel.Id)
			}
			conf.Ids = append(conf.Ids, chanIds...)
		}

	case config.ID_TYPE_USER:
		newIds := []string{}
		for _, id := range conf.Ids {
			resBody := []byte{}
			err := conf.discordRequest(http.MethodPost, "/users/@me/channels", strings.NewReader(fmt.Sprintf(`{"recipient_id":"%v"}`, id)), &resBody)
			if errors.Is(err, Err404) {
				fmt.Printf("Warning! %v is either not a user, or a dm cannot be opened with them. You have to manually get the id from them!\n", id)
			} else {
				tools.PanicIfErr(err)
			}
			resp := NewDMChannelResp{}
			err = json.Unmarshal(resBody, &resp)
			tools.PanicIfErr(err)
			newIds = append(newIds, resp.Id)
		}
		conf.Ids = newIds
	}

	if len(conf.Ids) == 1 {
		needChecking := 0
		if conf.Filter.MaxId != "" {
			maxMsgFetched, err := conf.FetchMsgId(conf.Ids[0], conf.Filter.MaxId)
			maxMsg = maxMsgFetched
			tools.PanicIfErr(err)
			needChecking++
			if conf.Filter.MinTime != 0 && conf.Filter.MinTime >= maxMsg.Timestamp {
				panic("BEFORE_ID is after BEFORE_TIME")
			}
		}
		if conf.Filter.MinId != "" {
			minMsgFetched, err := conf.FetchMsgId(conf.Ids[0], conf.Filter.MinId)
			tools.PanicIfErr(err)
			minMsg = minMsgFetched
			needChecking++
			if conf.Filter.MaxTime != 0 && minMsg.Timestamp >= conf.Filter.MaxTime {
				panic("AFTER_ID is after AFTER_TIME!")
			}
		}
		if needChecking == 2 {
			if minMsg.Timestamp > maxMsg.Timestamp {
				panic("Your BEFORE_ID is before AFTER_ID!")
			}
		}
	} else {
		if conf.Filter.MaxId != "" || conf.Filter.MinId != "" {
			fmt.Println("Alert! There is a Max ID or a Min ID, but multiple channels. This is not supported. If you want between certian dates on all channels, use BEFORE_TIME or AFER_TIME")
		}
	}

	if maxTime > maxMsg.Timestamp && maxMsg.Timestamp != 0 {
		maxTime = maxMsg.Timestamp
	}

	if minTime < minMsg.Timestamp {
		minTime = minMsg.Timestamp
	}

	for _, channel := range conf.Ids {
		if NumLeft == 0 {
			break
		}

		outputDir := tools.ParseTemplate(conf.ExportLocation, map[string]string{
			"CHANNEL_ID": channel,
		})

		os.Mkdir(outputDir, 0755)

		mediaDir := ""

		if conf.DownloadMedia {
			mediaDir = filepath.Join(outputDir, "media")
			err := os.Mkdir(mediaDir, 0755)
			if !os.IsExist(err) && err != nil {
				panic(err)
			}
		}

		file, err := os.Create(filepath.Join(outputDir, "content"+ext))

		tools.PanicIfErr(err)

		defer file.Close()

		lastMsgId := ""
		limit := 100

		if conf.UseLimit50 {
			limit = 50
		}

		theme := components.Theme{}

		if conf.ExportType == config.EXPORT_TYPE_HTML {
			theme.LoadTheme(conf.ExportHtmlThemeName)
			err = tools.CopyFile(filepath.Join(theme.BaseCss, "css", "style.css"), filepath.Join(outputDir, "base.css"))
			tools.PanicIfErr(err)
			err = tools.CopyFile(filepath.Join(theme.ThemeDir, "css", "style.css"), filepath.Join(outputDir, "style.css"))
			tools.PanicIfErr(err)

			err = os.Mkdir(filepath.Join(outputDir, "assets"), 0755)
			if os.IsExist(err) {

			} else {
				tools.PanicIfErr(err)
			}

			assetsBase, err := ioutil.ReadDir(filepath.Join(theme.BaseCss, "assets"))

			if os.IsNotExist(err) {
			} else {
				tools.PanicIfErr(err)
			}

			for _, asset := range assetsBase {
				name := asset.Name()
				err = tools.CopyFile(filepath.Join(theme.BaseCss, "assets", name), filepath.Join(outputDir, "assets", name))
				tools.PanicIfErr(err)
			}

			assetsTheme, err := ioutil.ReadDir(filepath.Join(theme.ThemeDir, "assets"))

			if os.IsNotExist(err) {
			} else {
				tools.PanicIfErr(err)
			}

			for _, asset := range assetsTheme {
				if asset.IsDir() {
					panic("The theme has dir assets. This is dis-allowed.")
				}
				name := asset.Name()
				err = tools.CopyFile(filepath.Join(theme.ThemeDir, "assets", name), filepath.Join(outputDir, "assets", name))
				tools.PanicIfErr(err)
			}

			resp := []byte{}
			conf.discordFetch(`/channels/`+channel, &resp)
			channelParsed := discord.Channel{}
			err = json.Unmarshal(resp, &channelParsed)
			tools.PanicIfErr(err)
			curTitle := ""

			switch channelParsed.Type {
			case discord.CHANNEL_TYPE_DM:
				curTitle = channelParsed.Recipients[0].Name
			default:
				curTitle = channelParsed.Name
			}
			file.WriteString(`<!DOCTYPE html><html>` + theme.HTMLHead(curTitle) + `<body>`)
			file.WriteString(theme.TopBar(curTitle, channelParsed.Type))

			if channelParsed.Type == discord.CHANNEL_TYPE_DM {
				file.WriteString(theme.StartDM(channelParsed.Recipients[0]))
			} // TODO: Add other channel types
		}

		for {
			fin := false
			allMsgs := conf.FetchChannelMessages(channel, lastMsgId, limit)

			if len(allMsgs) != limit {
				fin = true // don't break because you still need proccessing
			}

			for i, j := 0, len(allMsgs)-1; i < j; i, j = i+1, j-1 {
				allMsgs[i], allMsgs[j] = allMsgs[j], allMsgs[i]
			} //shameless stealing from so https://stackoverflow.com/questions/19239449/how-do-i-reverse-an-array-in-go

			prevMsg := discord.Message{}

			for _, msg := range allMsgs {
				for _, embed := range msg.Embeds {
					if embed.Type == discord.EMBED_IMAGE {
						msg.Attachments = append(msg.Attachments, discord.Attachment{
							ID:          "",
							Name:        "",
							Size:        0,
							Url:         embed.Url,
							Width:       embed.Thumbnail.Width,
							Height:      embed.Thumbnail.Height,
							ContentType: "image/webp",
						})
						if msg.Content == embed.Url {
							msg.Content = ""
						}
					}
				}

				attachments := ""
				if len(msg.Attachments) != 0 && conf.DownloadMedia {
					for _, attach := range msg.Attachments {
						attachments += fmt.Sprintf(`"%v",`, attach.Url)
						DownloadMedia(mediaDir, attach.Url, attach.Name)
					}
					attachments = attachments[:len(attachments)-1]
				}

				stickers := ""

				if len(msg.Stickers) != 0 && conf.DownloadMedia {
					for _, sticker := range msg.Stickers {
						stickers += fmt.Sprintf(`"%v",`, sticker.ID)
						DownloadMedia(mediaDir, sticker.URL(512), "STICKER_" + sticker.ID)
					}
					stickers = stickers[:len(stickers)-1]
				}

				switch conf.ExportType {
				case config.EXPORT_TYPE_TEXT:
					file.WriteString(tools.ParseTemplate(conf.ExportTextFormat+"\n", map[string]string{
						"AUTHOR_NAME":    msg.Author.Name,
						"AUTHOR_ID":      msg.Author.ID,
						"TIMESTAMP":      fmt.Sprint(msg.Timestamp),
						"WAS_EDITED":     fmt.Sprint(msg.IsEdited),
						"CONTENT":        msg.Content,
						"HAS_ATTACHMENT": fmt.Sprint(len(msg.Attachments) != 0),
						"ATTACHMENT_URL": attachments,
						"IS_REPLY":       fmt.Sprint(msg.IsReply),
						"IS_STICKER":     fmt.Sprint(msg.IsSticker),
						"STICKER_IDS":    stickers,
					}))
				case config.EXPORT_TYPE_HTML:
					msgTimestamp := discord.TimestampToTime(msg.Timestamp)
					sameDate := tools.SameDate(msgTimestamp, discord.TimestampToTime(prevMsg.Timestamp))

					if !sameDate {
						file.WriteString(theme.DateSeperator(msgTimestamp))
					}

					if len(msg.Attachments) != 0 {
						fmt.Printf("%#v\n", msg.Attachments)
					}

					file.WriteString(theme.MessageComponent(msg, prevMsg, prevMsg.Author.ID != msg.Author.ID || !sameDate || msg.IsReply))

					prevMsg = msg
					fin = true

				case config.EXPORT_TYPE_JSON:
					attachments := []JSONMetaAttachment{}
					newChanInfo := ParsedMessages[channel]
					for _, attach := range msg.Attachments {
						attachments = append(attachments, JSONMetaAttachment{
							Attachment: attach,
							AuthorID:   msg.Author.ID,
						})
						newChanInfo.AuthorAttachment[msg.Author.ID] = append(newChanInfo.AuthorAttachment[msg.Author.ID]) //TODO: newChanInfo.length??? I have no clue honestly. Im so fucking tired
					}

					// TODO: The rest of the things as well
					newChanInfo.Attachments = append(ParsedMessages[channel].Attachments, attachments...)
					newChanInfo.Msgs = append(newChanInfo.Msgs, msg)

					ParsedMessages[channel] = newChanInfo
				}
			}

			if fin {
				break
			}
		}

		if conf.ExportType == config.EXPORT_TYPE_HTML {
			file.WriteString(theme.MSG_INP_BAR)
			file.WriteString(`</body></html>`)
		}
	}
}

func (conf ConfigType) FetchChannelMessages(channel string, before string, limit int) []discord.Message {
	if len(before) != 0 {
		before = "&before=" + before
	}
	resp := []byte{}
	err := conf.discordFetch(fmt.Sprintf("/channels/%v/messages?limit=%v%v", channel, limit, before), &resp)
	tools.PanicIfErr(err)
	allMsgs := []discord.Message{}
	json.Unmarshal(resp, &allMsgs)
	return allMsgs
}
