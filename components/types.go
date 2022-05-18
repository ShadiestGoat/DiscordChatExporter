package components

type Theme struct {
	BaseCss,
	ThemeDir,
	MSG_WITH_PFP,
	MSG,
	MSG_INP_BAR,
	SVG_DM,
	SVG_CHANNEL,
	SVG_GROUP_DM,
	DATE_SEPERATOR,
	DM_START,
	HTML_HEAD,
	TOP_BAR,
	START_CHAN,
	REPLY,
	STICKER,
	GIF,
	OTHER_MIME,
	MENTION,
	CUSTOM_EMOJI,
	EMOJI_WRAPPER,
	REACTION,
	REACTION_WRAPPER,
	IMG string
	DownloadMedia bool
}

type SystemMsgs struct {
	USR_ADD,
	USR_RM,
	CALL,
	RENAME,
	ICON_CH,
	PIN,
	MEM_JOIN string
}
