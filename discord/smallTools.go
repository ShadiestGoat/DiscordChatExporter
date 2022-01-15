package discord

import (
	"fmt"
	"time"
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