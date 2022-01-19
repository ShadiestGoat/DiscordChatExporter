package config

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
	"github.com/joho/godotenv"
)

const MAX_ID = 999999999999999999

var localeReg = regexp.MustCompile(`.._..`)
var versionReg = regexp.MustCompile(`\d+\.\d+\.\d+`)
var quotesReg = regexp.MustCompile(`"`)

func parseMap(envMap map[string]envOpt, config *Config) {
	for key, opts := range envMap {
		val := os.Getenv(key)

		if len(val) == 0 {
			// no input
			if opts.Skip != nil {
				if opts.Skip(*config) {
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
			Skip: func(config Config) bool {
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

	if doAuto {
		superPropsOS := ""
		locale := ""
		discordPath, err := os.UserHomeDir()
		tools.PanicIfErr(err)

		switch runtime.GOOS {
		case "linux":
			superPropsOS = "Linux"
			locale = localeReg.FindString(os.Getenv("LANG"))
			if len(locale) == 0 {
				fmt.Println("Warning! Locale cannot be found! This *may* raise suspicion from discord!")
			} else {
				locale = locale[:2] + "-" + locale[3:]
			}
			discordPath += `/.config`
		case "darwin":
			superPropsOS = "Mac OS X"
			locale = localeReg.FindString(os.Getenv("LANG"))
			if len(locale) == 0 {
				fmt.Println("Warning! Locale cannot be found! This *may* raise suspicion from discord!")
			} else {
				locale = locale[:2] + "-" + locale[3:]
			}
			discordPath += `/Library/Application Support`
		case "windows":
			superPropsOS = "Windows"
			discordPath += `\AppData\Roaming`
			
			cmdLoc := exec.Command(`Get-Culture | select -exp Name`)
			localeInp, err := cmdLoc.Output()
			tools.PanicIfErr(err)
			locale = string(localeInp)
			panic(locale)
		default:
			panic("Auto header masking is not supported on your os!")
		}

		dExist := true
		normalD, err := ioutil.ReadDir(filepath.Join(discordPath, "discord"))
		
		if os.IsNotExist(err) {
			dExist = false
		} else {
			tools.PanicIfErr(err)
		}

		dCExist := true
		normalCD, err := ioutil.ReadDir(filepath.Join(discordPath, "discordcanary"))
		
		if os.IsNotExist(err) {
			dCExist = false
		} else {
			tools.PanicIfErr(err)
		}

		if dCExist && dExist && len(os.Getenv("USE_CANARY")) == 0 {
			panic("Both canary and stable discord detected. We don't know which one to pull from, use 'USE_CANARY'!")
		}

		if dExist && !dCExist && config.HeadersMask.UseCanary && len(os.Getenv("USE_CANARY")) != 0 {
			fmt.Println("Warning! Stable discord found but canary not found, and the preferance is for canary! We will be using false for USE_CANARY for this download")
			config.HeadersMask.UseCanary = false
		} else if !config.HeadersMask.UseCanary && dCExist && !dExist && len(os.Getenv("USE_CANARY")) != 0 {
			fmt.Println("Warning! Canary discord found but stable not found, and the preferance is for stable! We will be using true for USE_CANARY for this download")
			config.HeadersMask.UseCanary = true
		}

		if config.HeadersMask.UseCanary {
			config.HeadersMask.DomainPrefix = "canary."
		}

		discordVersion := ""

		discordFiles := normalD

		if config.HeadersMask.UseCanary {
			discordFiles = normalCD
		}

		for _, fileInfo := range discordFiles {
			if !fileInfo.IsDir() {
				continue
			}
			name := fileInfo.Name()
			if versionReg.MatchString(name) {
				discordVersion = name
			}
		}

		superPropsReleaseChannel := "stable"

		if config.HeadersMask.UseCanary {
			superPropsReleaseChannel = "canary"
		}
		
		config.HeadersMask.Locale = locale

		superInfo := ""

		if runtime.GOOS == "linux" {
			winMgr := os.Getenv("XDG_CURRENT_DESKTOP")
			if len(winMgr) == 0 {
				fmt.Println("Warning! Cannot find window manager!")
			} else {
				winMgr += ","
			}

			cmd := exec.Command("uname", "-r")
			osVersion, err := cmd.Output()
			tools.PanicIfErr(err)
			
			osInfo, err := ioutil.ReadFile("/etc/os-release")
			tools.PanicIfErr(err)
			distro := ""

			for _, line := range strings.Split(string(osInfo), "\n") {
				if line[:4] == "NAME" {
					distro = line[5:]
					distro = quotesReg.ReplaceAllString(distro, `\"`)
				}
			}
			
			superInfo = fmt.Sprintf(
				`{"os":"Linux","browser":"Discord Client","release_channel":"%v","client_version":"%v","os_version":"%v","os_arch":"x64","system_locale":"%v","window_manager":"%vunknown","distro":"%v","client_build_number":TODO:,"client_event_source":null}`, 
				superPropsReleaseChannel,
				discordVersion,
				osVersion,
				locale,
				winMgr,
				distro,
				)
		} else {
			superInfo = fmt.Sprintf(
				`{"os":"%v","browser":"Discord Client","release_channel":"%v","client_version":"%v","os_version":"TODO:","os_arch":"x64","system_locale":"%v","client_build_number":TODO:,"client_event_source":null}`,
				superPropsOS,
				discordVersion,
				superPropsReleaseChannel,
				locale,
			)
		}

		config.HeadersMask.SuperProperties = base64.StdEncoding.EncodeToString([]byte(superInfo))
	
		panic("Not implemented! So far: " + config.HeadersMask.SuperProperties)		
		
		// My props are these, but these need to be tested on mac & windows! 
		
		// {"os":"Linux","browser":"Discord Client","release_channel":"canary","client_version":"0.0.132","os_version":"5.15.12-arch1-1","os_arch":"x64","system_locale":"en-US","window_manager":"i3,unknown","distro":"\"Arch Linux\"","client_build_number":111095,"client_event_source":null}
		// {"os":"Mac OS X","browser":"Discord Client","release_channel":"stable","client_version":"0.0.264","os_version":"16.7.0","os_arch":"x64","system_locale":"en-US","client_build_number":110451,"client_event_source":null} this is a mac version
		// {"os":"Windows","browser":"Discord Client","release_channel":"stable","client_version":"1.0.9003","os_version":"10.0.22000","os_arch":"x64","system_locale":"en-GB","client_build_number":110451,"client_event_source":null}
		// client_version = \d+\.\d+\.\d+ in info dir :)
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