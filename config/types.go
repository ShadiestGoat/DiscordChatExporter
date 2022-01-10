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
	UseLimit50 bool
	Filter MsgFilter
}