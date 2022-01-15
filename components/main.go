package components

import (
	"fmt"
	"io/ioutil"
	// "io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/discord"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

func preParseHTML(html string) string {
	return extraWhitespaceReg.ReplaceAllString(tabReg.ReplaceAllString(newLineReg.ReplaceAllString(html, ""), " "), " ")
}

func (theme Theme) parseComponent(component string) string {
	comp, err := ioutil.ReadFile(filepath.Join(theme.ThemeDir, "components", component + ".html"))

	if os.IsNotExist(err) {
		baseComp, err := ioutil.ReadFile(filepath.Join(theme.BaseCss, "components", component + ".html"))
		tools.PanicIfErr(err)
		comp = baseComp
	} else {
		tools.PanicIfErr(err)
	}
	return preParseHTML(string(comp))
}

func (theme *Theme) LoadTheme(themeName string) {
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
	theme.START_DM = theme.parseComponent("startDm")
}

var newLineReg = regexp.MustCompile(`\n`)
var extraWhitespaceReg = regexp.MustCompile(`\s{2,}`)
var tabReg = regexp.MustCompile(`\t`)

const IMG = ``


func (theme Theme) MessageComponent(msg discord.Message, previousMsg discord.Message, firstMsg bool) string {
	content := msg.Content
	content = newLineReg.ReplaceAllString(content, "<br />")

	attachContent := ""

	for _, attach := range msg.Attachments {
		if attach.ContentType[:5] == "image" {
			attachContent += tools.ParseTemplate(IMG, map[string]string{
				"IMG_URL": attach.Url,
				"WIDTH": fmt.Sprint(0.8*float64(attach.Width)),
				"HEIGHT": fmt.Sprint(0.8*float64(attach.Height)),
			})
		} else {
			panic(attach.ContentType)
		}
	}

	if firstMsg {
		return tools.ParseTemplate(theme.MSG_WITH_PFP, map[string]string{
			"PFP": msg.Author.AvatarUrl(256),
			"USERNAME": msg.Author.Name,
			"DATE": discord.TimestampToTime(msg.Timestamp).Format("Mon 02/01/2006 03:04:05 PM"),
			"CONTENT": content,
			"ATTACH_CONTENT": attachContent,
			"ID": msg.ID,
			"REPLY_CONTENT": "",
		})
	} else {
		return tools.ParseTemplate(theme.MSG, map[string]string{
			"DATE": discord.TimestampToTime(msg.Timestamp).Format("15:04:05"),
			"CONTENT": content,
			"ATTACH_CONTENT": attachContent,
			"ID": msg.ID,
		})
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
		"ICON": icon,
		"TITLE": title,
	})
}

func (theme Theme) StartDM(author discord.Author) string {
	return tools.ParseTemplate(theme.START_DM, map[string]string{
		"TITLE": author.Name,
		"PFP": author.AvatarUrl(512),
	})
}