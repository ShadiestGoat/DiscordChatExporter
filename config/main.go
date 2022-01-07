package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const MAX_ID = 999999999999999999

func Load() Config {
	gottenIgnore := os.Getenv("IGNORE_ENV_FILE")
	ignoreEnv := false
	if len(gottenIgnore) != 0 {
		parsedIgnore, err := ParseBool(gottenIgnore)
		if err != nil {
			panic(err)
		}
		ignoreEnv = parsedIgnore
	}
	envFileName := os.Getenv("ENV_FILENAME")
	if len(envFileName) == 0 {
		envFileName = ".env"
	}
	if !ignoreEnv {
		godotenv.Load(envFileName)
	}

	config := Config{}

	envToken := os.Getenv("TOKEN")

	if len(envToken) == 0 {
		panic("TOKEN is required!")
	} else if envToken[:4] != "mfa." {
		panic("TOKEN is not a user token!")
	}
	
	config.Token = envToken

	IdType := os.Getenv("ID_TYPE")

	switch IdType {
		case "USER":
			config.IdType = ID_TYPE_USER
		case "CHANNEL":
			config.IdType = ID_TYPE_CHANNEL
		case "GUILD":
			config.IdType = ID_TYPE_GUILD
		default:
			if len(IdType) == 0 {
				config.IdType = ID_TYPE_CHANNEL
			} else {
				panic("ID_TYPE must be one of the allowed values!")
			}
	}

	idGotten := strings.Split(os.Getenv("ID"), " ")
	ids := []string{}
	
	for _, ogId := range idGotten {
		id, idOk := VerifySnowflake(ogId)
		if !idOk {
			fmt.Printf("WARNING! %v is not a valid discord id!", ogId)
		} else {
			ids = append(ids, id)
		}
	}

	config.Ids = ids

	envDwMedia := os.Getenv("DOWNLOAD_MEDIA")
	if len(envDwMedia) != 0 {
		parsedDwMedia, err := ParseBool(envDwMedia)
		if err != nil {
			panic(err)
		}
		config.DownloadMedia = parsedDwMedia
	}

	exportTypeGotten := os.Getenv("EXPORT_TYPE")

	switch exportTypeGotten {
		case "TEXT":
			config.ExportType = EXPORT_TYPE_TEXT
			envFormat := os.Getenv("EXPORT_PLAIN_FORMAT")
			if len(envFormat) == 0 {
				config.ExportTextFormat = `[{{$ID}}]: "{{$CONTENT}}"`
			} else {
				config.ExportTextFormat = envFormat
			}
		case "JSON":
			config.ExportType = EXPORT_TYPE_JSON
			envExportJson := os.Getenv("EXPORT_JSON_TOOLS")
			if len(envExportJson) == 0 {
				config.ExportJsonMeta = true
			} else {
				parsedJsonMeta, err := ParseBool(envExportJson)
				if err != nil {
					panic(err)
				}
				config.ExportJsonMeta = parsedJsonMeta
			}
		case "HTML":
			config.ExportType = EXPORT_TYPE_HTML
			config.ExportHtmlThemeName = "black"
			envThemeName := os.Getenv("EXPORT_HTML_THEME")
			if len(envThemeName) != 0 {
				config.ExportHtmlThemeName = envThemeName
			}
		default:
			if len(exportTypeGotten) == 0 {
				config.IdType = ID_TYPE_CHANNEL
			} else {
				panic("EXPORT_TYPE must be one of the allowed values!")
			}
	}

	envExportLoc := os.Getenv("EXPORT_LOCATION")
	if len(envExportLoc) == 0 {
		envExportLoc = filepath.Join("output", "{{$ID}}")
	}
	config.ExportLocation = envExportLoc

	envNumLimit := os.Getenv("MSG_LIMIT_NUM")

	if len(envNumLimit) != 0 && envNumLimit != "all" {
		parsedNum, err := strconv.ParseInt(envNumLimit, 10, 64)
		if err != nil {
			panic(err)
		}
		config.Filter.NumMax = int(parsedNum)
	} else {
		config.Filter.NumMax = MAX_ID
	}
	if config.Filter.NumMax <= 0 {
		panic("Number limit may not be below 1!")
	}

	envBeforeId := os.Getenv("BEFORE_ID")
	if len(envBeforeId) != 0 {
		config.Filter.MaxId = envBeforeId
	}
	envAfterId := os.Getenv("AFTER_ID")
	if len(envAfterId) != 0 {
		config.Filter.MinId = envAfterId
	}
	envBeforeTime := os.Getenv("BEFORE_TIME")
	parsedMaximum := 0
	if len(envBeforeTime) != 0 {
		parsedMax, err := strconv.ParseInt(envBeforeTime, 10, 64)
		if err != nil {
			panic(err)
		}
		config.Filter.MaxTime = int(time.Unix(parsedMax, 0).UnixMicro())
		parsedMaximum = int(parsedMax)
	}
	envAfterTime := os.Getenv("AFTER_TIME")
	if len(envAfterTime) != 0 {
		parsedMin, err := strconv.ParseInt(envAfterTime, 10, 64)
		if err != nil {
			panic(err)
		}
		config.Filter.MaxTime = int(time.Unix(parsedMin, 0).UnixMicro())
		if len(envBeforeTime) != 0 && parsedMaximum < int(parsedMin) {
			panic("AFTER_TIME is after BEFORE_TIME")
		}
	}
	return config
}