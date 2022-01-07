package fetcher

import "fmt"

func (author Author) AvatarUrl() string {
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/376079696489742338/%v.webp?size=512", author.Avatar)
}