package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
	"github.com/joho/godotenv"
)

const MAX_ID = 999999999999999999

func Load() Config {
	config := Config{}

	gottenIgnore := os.Getenv("IGNORE_ENV_FILE")
	ignoreEnv := false
	if len(gottenIgnore) != 0 {
		parsedIgnore, err := ParseBool(gottenIgnore)
		tools.PanicIfErr(err)
		ignoreEnv = parsedIgnore
	}
	
	envFileName := os.Getenv("ENV_FILENAME")

	if len(envFileName) == 0 {
		envFileName = ".env"
	}

	if !ignoreEnv {
		godotenv.Load(envFileName)
	}

	doAuto := false
	idType := ""
	exportType := ""
	ids := ""
	limNum := ""

	maxTime := ""
	minTime := ""

	envMap := map[string]envOpt{
		"HM_AUTO": {
			DefaultBool: true,
			Type: ENV_TYPE_BOOL,
			PointBool: &doAuto,
		},
		"HM_USER_AGENT": {
			Type: ENV_TYPE_STRING,
			PointString: &config.HeadersMask.UserAgent,
			Skip: func(config Config) bool {
				return doAuto
			},
			NoDefault: true,
		},
		"USE_CANARY": {
			Type: ENV_TYPE_BOOL,
			PointBool: &config.HeadersMask.UseCanary,
			Skip: func(config Config) bool {
				return doAuto
			},
			NoDefault: true,
		},
		"HM_LOCALE": {
			Type: ENV_TYPE_STRING,
			PointString: &config.HeadersMask.Locale,
			Skip: func(config Config) bool {
				return doAuto
			},
			NoDefault: true,
		},
		"HM_SUPER_PROPS": {
			Type: ENV_TYPE_STRING,
			PointString: &config.HeadersMask.SuperProperties,
			Skip: func(config Config) bool {
				return doAuto
			},
			NoDefault: true,
		},
		"TOKEN": {
			Type: ENV_TYPE_STRING,
			NoDefault: true,
			PointString: &config.Token,
		},
		"ID_TYPE": {
			Type: ENV_TYPE_STRING,
			DefaultString: "CHANNEL",
			PointString: &idType,
		},
		"IGNORE_NSFW_CHANNEL": {
			Type: ENV_TYPE_BOOL,
			PointBool: &config.IgnoreNsfw,
			DefaultBool: false,
		},
		"ID": {
			Type: ENV_TYPE_STRING,
			PointString: &ids,
			NoDefault: true,
		},
		"DOWNLOAD_MEDIA": {
			DefaultBool: true,
			Type: ENV_TYPE_BOOL,
			PointBool: &config.DownloadMedia,
		},
		"EXPORT_TYPE": {
			Type: ENV_TYPE_STRING,
			PointString: &exportType,
			DefaultString: "JSON",
		},
		"EXPORT_JSON_TOOLS": {
			Type: ENV_TYPE_BOOL,
			DefaultBool: true,
			PointBool: &config.ExportJsonMeta,
		},
		"EXPORT_PLAIN_FORMAT": {
			Type: ENV_TYPE_STRING,
			DefaultString: `[{{%CHANNEL_ID}}]: "{{%CONTENT}}"`,
			PointString: &config.ExportTextFormat,
		},
		"EXPORT_HTML_THEME": {
			Type: ENV_TYPE_STRING,
			DefaultString: "dark",
			PointString: &config.ExportHtmlThemeName,
		},
		"EXPORT_LOCATION": {
			Type: ENV_TYPE_STRING,
			DefaultString: filepath.Join("output", "{{%CHANNEL_ID}}"),
			PointString: &config.ExportLocation,
		},
		"MSG_LIMIT_NUM": {
			Type: ENV_TYPE_STRING,
			DefaultString: "all",
			PointString: &limNum,
		},
		"BEFORE_ID": {
			Type: ENV_TYPE_STRING,
			PointString: &config.Filter.MaxId,
		},
		"AFTER_ID": {
			Type: ENV_TYPE_STRING,
			PointString: &config.Filter.MinId,
		},
		"AFTER_TIME": {
			Type: ENV_TYPE_STRING,
			PointString: &minTime,
		},
		"BEFORE_TIME": {
			Type: ENV_TYPE_STRING,
			PointString: &maxTime,
		},
		"USE_LIMIT_50": {
			Type: ENV_TYPE_BOOL,
			DefaultBool: true,
			PointBool: &config.UseLimit50,
		},
	}

	for key, opts := range envMap {
		val := os.Getenv(key)

		if len(val) == 0 {
			// no input
			if opts.Skip(config) {
				continue
			}
			if opts.NoDefault {
				panic(fmt.Sprintf("%v was not found, but is needed for startup! Check your env file!", key))
			}
			switch opts.Type {
			case ENV_TYPE_BOOL:
				*opts.PointBool = opts.DefaultBool
			case ENV_TYPE_STRING:
				*opts.PointString = opts.DefaultString
			}
		} else {
			switch opts.Type {
			case ENV_TYPE_BOOL:
				parsedBoolP, err := ParseBool(val)
				tools.PanicIfErr(err)
				*opts.PointBool = parsedBoolP
			case ENV_TYPE_STRING:
				*opts.PointString = val
			}
		}
	}

	switch idType {
	case "USER":
		config.IdType = ID_TYPE_USER
	case "CHANNEL":
		config.IdType = ID_TYPE_CHANNEL
	case "GUILD":
		config.IdType = ID_TYPE_GUILD
	default:
		panic("ID_TYPE must be one of the allowed values!")
	}

	idGotten := strings.Split(ids, " ")
	
	for _, ogId := range idGotten {
		id, idOk := VerifySnowflake(ogId)
		if !idOk {
			fmt.Printf("WARNING! %v is not a valid discord id!\n", ogId)
		} else {
			idGotten = append(idGotten, id)
		}
	}

	if len(idGotten) == 0 {
		panic("No valid IDs were found!")
	}

	config.Ids = idGotten

	switch exportType {
	case "HTML":
		config.ExportType = EXPORT_TYPE_HTML
	case "TEXT":
		config.ExportType = EXPORT_TYPE_TEXT
	case "JSON":
		config.ExportType = EXPORT_TYPE_JSON
	}

	if len(limNum) != 0 && limNum != "all" {
		parsedNum, err := strconv.ParseInt(limNum, 10, 64)
		tools.PanicIfErr(err)
		config.Filter.NumMax = int(parsedNum)
	} else {
		config.Filter.NumMax = MAX_ID
	}

	if config.Filter.NumMax < 1 {
		panic("Number limit may not be below 1!")
	}

	if len(maxTime) != 0 {
		parsedMax, err := strconv.ParseInt(maxTime, 10, 64)
		tools.PanicIfErr(err)
		config.Filter.MaxTime = int(time.Unix(parsedMax, 0).UnixMicro())
	}

	if len(minTime) != 0 {
		parsedMin, err := strconv.ParseInt(minTime, 10, 64)
		tools.PanicIfErr(err)
		config.Filter.MaxTime = int(time.Unix(parsedMin, 0).UnixMicro())
		if config.Filter.MaxTime < config.Filter.MinTime {
			panic("AFTER_TIME is after BEFORE_TIME")
		}
	}

	if doAuto {
		panic("Not implemented!")
		// superInfo := fmt.Sprintf(`{"os":"Linux","browser":"Discord Client","release_channel":"canary","client_version":"0.0.132","os_version":"5.15.12-arch1-1","os_arch":"x64","system_locale":"en-US","window_manager":"i3,unknown","distro":"\"Arch Linux\"","client_build_number":111095,"client_event_source":null}`)
		
		// My props are these, but these need to be tested on mac & windows! 
		
		// {"os":"Linux","browser":"Discord Client","release_channel":"canary","client_version":"0.0.132","os_version":"5.15.12-arch1-1","os_arch":"x64","system_locale":"en-US","window_manager":"i3,unknown","distro":"\"Arch Linux\"","client_build_number":111095,"client_event_source":null}
		// {"os":"Mac OS X","browser":"Discord Client","release_channel":"stable","client_version":"0.0.264","os_version":"16.7.0","os_arch":"x64","system_locale":"en-US","client_build_number":110451,"client_event_source":null} this is a mac version
		// {"os":"Windows","browser":"Discord Client","release_channel":"stable","client_version":"1.0.9003","os_version":"10.0.22000","os_arch":"x64","system_locale":"en-GB","client_build_number":110451,"client_event_source":null}
		// client_version = \d+\.\d+\.\d+ in info dir :)
		// if canary release_channel = canary, else its stable
		// browser = Discord Client
		// OSes = Linux, Windows, Mac OS X,
		// os_version
		// os_arch
		// system_locale on win its Get-Culture | select -exp Name, on unix its just LANG os env var and then parse using .._.. and replace the _ w/ '-'
		// client_build_number ??? idk
		// client_event_source = null
		// LINUX:
		// window_manager = `$XDG_CURRENT_DESKTOP,unknown`
		// read /etc/os-release and the NAME

		// Fetch canary or discord. First check if both exist. If both exist, check which one is prefered, if no pref input panic. If the prefered doesn't exist do a warning. 
		// 
		// TODO: Automatically fetch the neededed items, specifically in Config.HeadingMask -
	}

	if config.HeadersMask.UseCanary {
		config.HeadersMask.DomainPrefix = "canary."
	}

	return config
}