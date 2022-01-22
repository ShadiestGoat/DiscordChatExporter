package config

type IdType int8

const (
	ID_TYPE_USER IdType = iota
	ID_TYPE_CHANNEL
	ID_TYPE_GUILD
)

type ExportType int8

const (
	EXPORT_TYPE_TEXT ExportType = iota
	EXPORT_TYPE_JSON
	EXPORT_TYPE_HTML
)

type MsgFilter struct {
	NumMax int
	MinId string
	MaxId string
	MinTime int
	MaxTime int
}

type Config struct {
	Token string
	IdType IdType
	Ids []string
	ExportType ExportType
	DownloadMedia bool
	ExportLocation string
	ExportHtmlThemeName string
	ExportJsonMeta bool
	ExportTextFormat string
	IgnoreNsfw bool
	IgnoreSystemMsgs bool
	UseLimit50 bool
	Filter MsgFilter
	HeadersMask HeadersMask
}

type HeadersMask struct {
	UseCanary bool
	UserAgent string
	Locale string
	SuperProperties string
	DomainPrefix string
	DiscordVersion string
}

type envOpt struct {
	DefaultString string
	DefaultBool bool
	Type envType
	PointBool *bool
	PointString *string
	NoDefault bool
	Skip func(config Config, selfExist bool) bool
}

type envType int8

const (
	ENV_TYPE_STRING envType = iota
	ENV_TYPE_BOOL
)
