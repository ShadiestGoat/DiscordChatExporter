package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
	"github.com/joho/godotenv"
)

const MAX_ID = 999999999999999999

func parseMap(envMap map[string]envOpt, config *Config) {
	for key, opts := range envMap {
		val := os.Getenv(key)

		if len(val) == 0 {
			// no input
			if opts.Skip != nil {
				if opts.Skip(*config, false) {
					continue
				}
			}

			if opts.NoDefault {
				panic(fmt.Sprintf("'%v' was not found, but is needed for startup! Check your env file!", key))
			}
			switch opts.Type {
			case ENV_TYPE_BOOL:
				*opts.PointBool = opts.DefaultBool
			case ENV_TYPE_STRING:
				*opts.PointString = opts.DefaultString
			}
		} else {
			if opts.Skip != nil {
				if opts.Skip(*config, true) {
					continue
				}
			}
			switch opts.Type {
			case ENV_TYPE_BOOL:
				parsedBoolP, err := ParseBool(val)
				if err != nil {
					// It only returns the bad bool error
					panic(fmt.Sprintf("'%v' has incorrect syntax for a boolean type. Please confirm the syntax with the wiki!", key))
				}
				*opts.PointBool = parsedBoolP
			case ENV_TYPE_STRING:
				*opts.PointString = val
			}
		}
	}
}

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

	priorityMap := map[string]envOpt{
		"HM_AUTO": {
			DefaultBool: true,
			Type: ENV_TYPE_BOOL,
			PointBool: &doAuto,
		},
	}

	envMap := map[string]envOpt{
		"HM_USER_AGENT": {
			Type: ENV_TYPE_STRING,
			PointString: &config.HeadersMask.UserAgent,
			Skip: func(config Config, selfExist bool) bool {
				if doAuto && selfExist {
					fmt.Println("Warning! HM_AUTO is set, so HM_USER_AGENT will not be used!")
				}
				return doAuto
			},
			NoDefault: true,
		},
		"USE_CANARY": {
			Type: ENV_TYPE_BOOL,
			PointBool: &config.HeadersMask.UseCanary,
		},
		"HM_LOCALE": {
			Type: ENV_TYPE_STRING,
			PointString: &config.HeadersMask.Locale,
			Skip: func(config Config, selfExist bool) bool {
				if doAuto && selfExist {
					fmt.Println("Warning! HM_AUTO is set, so HM_LOCALE will not be used!")
				}
				return doAuto
			},
			NoDefault: true,
		},
		"HM_DISCORD_VERSION": {
			Type: ENV_TYPE_STRING,
			PointString: &config.HeadersMask.DiscordVersion,
			Skip: func(config Config, selfExist bool) bool {
				if doAuto && selfExist {
					fmt.Println("Warning! HM_AUTO is set, so HM_DISCORD_VERSION will not be used!")
				}
				return doAuto
			},
			NoDefault: true,
		},
		"HM_SUPER_PROPS": {
			Type: ENV_TYPE_STRING,
			PointString: &config.HeadersMask.SuperProperties,
			Skip: func(config Config, selfExist bool) bool {
				if doAuto && selfExist {
					fmt.Println("Warning! HM_AUTO is set, so HM_SUPER_PROPS will not be used!")
				}
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

	parseMap(priorityMap, &config)
	parseMap(envMap, &config)

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

	if config.HeadersMask.UseCanary {
		config.HeadersMask.DomainPrefix = "canary."
	}

	if doAuto {
		switch runtime.GOOS {
			case "windows", "linux", "darwin":
			default:
				panic("Auto header masking is not supported on your os!")
		}

		config.HeadersMask.Auto()
	}


	return config
}