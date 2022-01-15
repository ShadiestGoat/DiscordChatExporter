package discord

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

type MsgType int8

const (
	MESSAGE_TYPE_NORMAL MsgType = iota

)

type RawMessage struct {
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
	Type MsgType `json:"type"`
	IsEdited bool
	IsReply bool
}

type Message struct {
	RawMessage
	ReplyTo RawMessage
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
	Id string `json:"id"`
	Nsfw bool `json:"nsfw"`
	Name string `json:"name"`
	Type ChannelType `json:"type"`
	Recipients []Author `json:"recipients"`
}

type Author struct {
	ID string `json:"id"`
	Name string `json:"username"`
	Avatar string `json:"avatar"`
	Discriminator string `json:"discriminator"`
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