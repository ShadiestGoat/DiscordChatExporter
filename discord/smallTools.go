package discord

import (
	"fmt"
	"time"
)

func (author Author) AvatarUrl(size int) string {
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%v.webp?size=%v", author.ID, author.Avatar, size)
}

func TimestampToTime(timestamp int) time.Time {
	return time.Unix(int64(timestamp/1000000), 0)
}