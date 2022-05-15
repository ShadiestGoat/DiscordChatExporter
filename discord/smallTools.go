package discord

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

func (author Author) URL(size int) string {
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%v.webp?size=%v", author.ID, author.Avatar, size)
}

func TimestampToTime(timestamp int) time.Time {
	return time.Unix(int64(timestamp/1000000), 0)
}

func (sticker Sticker) URL(size int) string {
	return fmt.Sprintf("https://media.discordapp.net/stickers/%v.webp?size=%v", sticker.ID, size)
}

// Thank you to bwmarrin/discordgo since I basically copied this from them :)
func IDToTimestamp(id string) int {
	i, err := strconv.ParseInt(id, 10, 64)
	tools.PanicIfErr(err)

	timestamp := (i >> 22) + 1420070400000

	return int(time.Unix(0, timestamp*1000000).UnixMicro())
}

// Reverse of IDToTimestamp, but note that because of the nature of snowflakes this will output the exact same output.
//
// The last 6 digits are different Can anyone confirm this? I'm pretty sure this is true, but idk man
func TimestampToID(timestamp int) string {
	id := (timestamp/1000 - 1420070400000) << 22

	return fmt.Sprint(id)
}

var badCharset = regexp.MustCompile(`; ?charset=.+?( |$)`)

func (attachment Attachment) MediaName() string {
	split := strings.Split(attachment.ContentType, "/")
	// remove ;charset=
	extSuf := strings.Split(split[1], ";")
	ext := "." + extSuf[0]
	badCharset.ReplaceAllString(ext, "")

	if ext == ".plain" {
		ext = filepath.Ext(attachment.Url)
	} else if ext == "svg+xml" {
		ext = "svg"
	}

	return attachment.ID + ext
}
