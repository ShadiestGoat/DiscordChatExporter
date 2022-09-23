package components

import (
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/discord"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

func preParseHTML(html string) string {
	return extraWhitespaceReg.ReplaceAllString(tabReg.ReplaceAllString(newLineReg.ReplaceAllString(html, ""), " "), " ")
}

func (theme Theme) parseComponent(component string) string {
	comp, err := os.ReadFile(filepath.Join(theme.ThemeDir, "components", component+".html"))

	if os.IsNotExist(err) {
		baseComp, err := os.ReadFile(filepath.Join(theme.BaseCss, "components", component+".html"))
		tools.PanicIfErr(err)
		comp = baseComp
	} else {
		tools.PanicIfErr(err)
	}
	return preParseHTML(string(comp))
}

func (theme *Theme) LoadTheme(themeName string, DW_MEDIA bool) {
	themesMainLoc, err := os.Stat("themes")
	if os.IsNotExist(err) || !themesMainLoc.IsDir() {
		panic("'themes' folder is not available! This folder must be present!")
	} else {
		tools.PanicIfErr(err)
	}

	themeDir := filepath.Join("themes", themeName)

	themeLoc, err := os.Stat(themeDir)

	if os.IsNotExist(err) || !themeLoc.IsDir() {
		panic("Cannot find theme!")
	} else {
		tools.PanicIfErr(err)
	}

	theme.ThemeDir = themeDir
	theme.BaseCss = filepath.Join("themes", "base")

	theme.SVG_CHANNEL = theme.parseComponent("SVGchan")
	theme.SVG_DM = theme.parseComponent("SVGdm")
	theme.SVG_GROUP_DM = theme.parseComponent("SVGgroup")
	theme.MSG = theme.parseComponent("normalMsg")
	theme.MSG_INP_BAR = theme.parseComponent("inputBar")
	theme.MSG_WITH_PFP = theme.parseComponent("msgStarter")
	theme.DATE_SEPERATOR = theme.parseComponent("dateSeperator")
	theme.HTML_HEAD = theme.parseComponent("htmlHead")
	theme.TOP_BAR = theme.parseComponent("topBar")
	theme.START_CHAN = theme.parseComponent("startChan")
	theme.IMG = theme.parseComponent("img")
	theme.STICKER = theme.parseComponent("sticker")
	theme.REPLY = theme.parseComponent("reply")
	theme.GIF = theme.parseComponent("gifs")
	theme.OTHER_MIME = theme.parseComponent("otherMime")
	theme.MENTION = theme.parseComponent("mention")
	theme.CUSTOM_EMOJI = theme.parseComponent("customEmoji")
	theme.EMOJI_WRAPPER = theme.parseComponent("emojiWrapper")
	theme.REACTION_WRAPPER = theme.parseComponent("reactionWrapper")
	theme.REACTION = theme.parseComponent("reaction")

	theme.DownloadMedia = DW_MEDIA
}

var newLineReg = regexp.MustCompile(`\n`)
var extraWhitespaceReg = regexp.MustCompile(`\s{2,}`)
var tabReg = regexp.MustCompile(`\t`)
var customEmojiReg = regexp.MustCompile(`&lt;a?:[^\s]+:\d+&gt;`)

func (theme Theme) MessageComponent(msg discord.Message, previousMsg discord.Message, firstMsg bool) string {
	content := msg.Content

	attachContent := ""

	content = template.HTMLEscapeString(content)

	content = newLineReg.ReplaceAllString(content, "<br />")
	content, _ = MDToHTML(content)

	tmpEmojiContent := customEmojiReg.ReplaceAllString(content, "")
	tmpEmojiContent = strings.TrimSpace(tmpEmojiContent)
	cssClassses := "text-emoji "
	if tmpEmojiContent == "" {
		cssClassses += "big"
	}

	for customEmojiReg.MatchString(content) {
		loc := customEmojiReg.FindStringIndex(content)
		emojiRaw := content[loc[0]:loc[1]]
		emojiId := emojiRaw[len(emojiRaw)-22 : len(emojiRaw)-4]
		format := "webp"
		if emojiRaw[5] == []byte("a")[0] {
			format = "gif"
		}
		// as used by discord
		href := discord.EmojiURL(emojiId, format)
		content = content[:loc[0]] + tools.ParseTemplate(theme.EMOJI_WRAPPER, map[string]string{
			"CONTENT": tools.ParseTemplate(theme.CUSTOM_EMOJI, map[string]string{
				"URL": href,
				"CSS_CLASSES": cssClassses,
			}),
		})  + content[loc[1]:]
	}

	for _, attach := range msg.Attachments {
		if attach.ContentType[:5] == "image" {
			mediaUrl := attach.Url
			if theme.DownloadMedia {
				mediaUrl = filepath.Join("media", attach.MediaName())
			}
			attachContent += tools.ParseTemplate(theme.IMG, map[string]string{
				"IMG_URL": mediaUrl,
				"WIDTH":   fmt.Sprint(math.Floor(0.8 * float64(attach.Width))),
				"HEIGHT":  fmt.Sprint(math.Floor(0.8 * float64(attach.Height))),
			})
		} else {
			attachContent += tools.ParseTemplate(theme.OTHER_MIME, map[string]string{
				"FILENAME": attach.Name,
			})
		}
	}

	gifContents := ""

	for _, embed := range msg.Embeds {
		if embed.Type == discord.EMBED_GIF {
			gifContents += tools.ParseTemplate(theme.GIF, map[string]string{
				"VIDEO_URL": embed.Url,
				"WIDTH":     fmt.Sprint(float64(embed.Thumbnail.Width) * 0.7),
				"HEIGHT":    fmt.Sprint(float64(embed.Thumbnail.Height) * 0.7),
			})
			if msg.Content == embed.GifContentUrl {
				content = ""
			}
		}
	}

	stickerContent := ""

	for _, sticker := range msg.Stickers {
		stickerContent += tools.ParseTemplate(theme.STICKER, map[string]string{
			"IMG_URL": sticker.URL(160),
		})
	}

	if msg.IsSystemType {
		return "TODO:"
	} else {
		reactions := ""
		for _, reaction := range msg.Reactions {
			emojiStr := ""
			if reaction.Emoji.ID == "" {
				emojiStr = reaction.Emoji.Name
			} else {
				emojiStr = tools.ParseTemplate(theme.CUSTOM_EMOJI, map[string]string{
					"CSS_CLASSES": "reaction-emoji",
					"URL": discord.EmojiURL(reaction.Emoji.ID, "webp"),
				})
			}
			reactions += tools.ParseTemplate(theme.REACTION, map[string]string{
				"EMOJI": emojiStr,
				"COUNT": fmt.Sprint(reaction.Count),
			})
		}
	
		reactions = tools.ParseTemplate(theme.REACTION_WRAPPER, map[string]string{
			"CONTENT": reactions,
		})
	
		for _, mention := range msg.Mentions {
			idReg := regexp.MustCompile(fmt.Sprintf(`&lt;@!?%v&gt;`, mention.ID))
			content = idReg.ReplaceAllString(content, tools.ParseTemplate(theme.MENTION, map[string]string{
				"MENTION_NAME":   mention.Name,
				"MENTION_PREFIX": "@",
			}))
		}

		if firstMsg {
			replyContent := ""

			if msg.IsReply {
				replyContent = tools.ParseTemplate(theme.REPLY, map[string]string{
					"PFP":     msg.ReplyTo.Author.URL(16),
					"NAME":    msg.ReplyTo.Author.Name,
					"CONTENT": msg.ReplyTo.Content,
				})
			}

			return tools.ParseTemplate(theme.MSG_WITH_PFP, map[string]string{
				"PFP":             msg.Author.URL(256),
				"USERNAME":        msg.Author.Name,
				"DATE":            discord.TimestampToTime(msg.Timestamp).Format("Mon 02/01/2006 03:04:05 PM"),
				"CONTENT":         content,
				"ATTACH_CONTENT":  attachContent,
				"ID":              msg.ID,
				"REPLY_CONTENT":   replyContent,
				"STICKER_CONTENT": stickerContent,
				"GIFS":            gifContents,
				"REACTIONS":	   reactions,
			})
		} else {
			return tools.ParseTemplate(theme.MSG, map[string]string{
				"DATE":            discord.TimestampToTime(msg.Timestamp).Format("15:04"),
				"CONTENT":         content,
				"ATTACH_CONTENT":  attachContent,
				"ID":              msg.ID,
				"STICKER_CONTENT": stickerContent,
				"GIFS":            gifContents,
				"REACTIONS":	   reactions,
			})
		}
	}
}

func (theme Theme) DateSeperator(date time.Time) string {
	return tools.ParseTemplate(theme.DATE_SEPERATOR, map[string]string{
		"DATE": date.Format("January 2, 2006"),
	})
}

func (theme Theme) HTMLHead(title string) string {
	return tools.ParseTemplate(theme.HTML_HEAD, map[string]string{
		"TITLE": title,
	})
}

func (theme Theme) TopBar(title string, channelType discord.ChannelType) string {
	icon := ""

	switch channelType {
	case discord.CHANNEL_TYPE_DM:
		icon = theme.SVG_DM
	case discord.CHANNEL_TYPE_GROUP_DM:
		icon = theme.SVG_GROUP_DM
	default:
		icon = theme.SVG_CHANNEL
	}

	return tools.ParseTemplate(theme.TOP_BAR, map[string]string{
		"ICON":  icon,
		"TITLE": title,
	})
}

func (theme Theme) StartDM(author discord.Author) string {
	return tools.ParseTemplate(theme.START_CHAN, map[string]string{
		"TITLE": author.Name,
		"ICON":  author.URL(512),
	})
}

func (theme Theme) StartChannel(channel discord.Channel) string {
	return tools.ParseTemplate(theme.START_CHAN, map[string]string{
		"TITLE": "Welcome to #" + channel.Name,
		"ICON":  "assets/channelStart.png",
	})
}

func (theme Theme) StartGroupDM(groupDM discord.Channel) string {
	return tools.ParseTemplate(theme.START_CHAN, map[string]string{
		"TITLE": groupDM.Name,
		"ICON":  groupDM.Icon,
	})
}
