package discord

type EmbedType int8

const (
	EMBED_LINK = iota
	EMBED_GIF
	EMBED_IGNORE
	EMBED_IMAGE
)

type Embed struct {
	Type          EmbedType
	Url           string
	Title         string
	Description   string
	Color         string
	Thumbnail     EmbedImageThumbnail
	GifContentUrl string
}

type EmbedVideo struct {
	Url string `json:"url"`
}

type MsgType int8

const (
	MSGT_DEFAULT MsgType = iota
	MSGT_RECIPIENT_ADD
	MSGT_RECIPIENT_REMOVE
	MSGT_CALL
	MSGT_CHANNEL_NAME_CHANGE
	MSGT_CHANNEL_ICON_CHANGE
	MSGT_CHANNEL_PINNED_MESSAGE
	MSGT_GUILD_MEMBER_JOIN
	MSGT_USER_PREMIUM_GUILD_SUBSCRIPTION
	MSGT_USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_1
	MSGT_USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_2
	MSGT_USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_3
	MSGT_CHANNEL_FOLLOW_ADD
	_
	MSGT_GUILD_DISCOVERY_DISQUALIFIED
	MSGT_GUILD_DISCOVERY_REQUALIFIED
	MSGT_GUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING
	MSGT_GUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING
	MSGT_THREAD_CREATED
	MSGT_REPLY
	MSGT_CHAT_INPUT_COMMAND
	MSGT_THREAD_STARTER_MESSAGE
	MSGT_GUILD_INVITE_REMINDER
	MSGT_CONTEXT_MENU_COMMAND
)

type EmbedImageThumbnail struct {
	Width  int
	Height int
	Url    string `json:"url"`
}

type Sticker struct {
	ID string `json:"id"`
}

type ReplyMsg struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Author  Author `json:"author"`
}

type Message struct {
	ReplyMsg
	Channel          string       `json:"channel_id"`
	Attachments      []Attachment `json:"attachments"`
	Embeds           []Embed      `json:"embeds"`
	Mentions         []Author     `json:"mentions"`
	MentionRoles     []string     `json:"mention_roles"`
	MentionsEveryone bool         `json:"mention_everyone"`
	Timestamp        int          `json:"timestamp"`
	Type             MsgType      `json:"type"`
	Stickers         []Sticker    `json:"sticker_items"`
	IsEdited         bool
	IsReply          bool
	IsSystemType     bool
	HasSticker       bool
	ReplyTo          ReplyMsg
}

type ChannelType int8

const (
	CHANNEL_TYPE_GUILD_TEXT ChannelType = iota
	CHANNEL_TYPE_DM
	CHANNEL_TYPE_GUILD_VOICE
	CHANNEL_TYPE_GROUP_DM
	CHANNEL_TYPE_GUILD_CATEGORY
	CHANNEL_TYPE_GUILD_NEWS
	CHANNEL_TYPE_GUILD_STORE
	CHANNEL_TYPE_GUILD_NEWS_THREAD
	CHANNEL_TYPE_GUILD_PUBLIC_THREAD
	CHANNEL_TYPE_GUILD_PRIVATE_THREAD
	CHANNEL_TYPE_GUILD_STAGE_VOICE
)

type Channel struct {
	Id         string      `json:"id"`
	Nsfw       bool        `json:"nsfw"`
	Name       string      `json:"name"`
	Type       ChannelType `json:"type"`
	Recipients []Author    `json:"recipients"`
	Icon       string      `json:"icon"`
}

type Author struct {
	ID            string `json:"id"`
	Name          string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}

type Attachment struct {
	ID          string `json:"id"`
	Name        string `json:"filename"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ContentType string `json:"content_type"`
}
