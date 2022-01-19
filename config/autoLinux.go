//go:build linux

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

func (mask *HeadersMask) Auto() {
	locale := localeReg.FindString(os.Getenv("LANG"))

	if len(locale) == 0 {
		fmt.Println("Warning! Locale cannot be found! This *may* raise suspicion from discord!")
	} else {
		locale = locale[:2] + "-" + locale[3:]
	}

	mask.Locale = locale
	mask.PullDiscordVers("~/.config")
	mask.UserAgent = fmt.Sprintf("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/%v Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36", mask.DiscordVersion)
	
	releaseChan := "stable"

	if mask.UseCanary {
		releaseChan = "canary"
	}
	
	cmd := exec.Command("uname", "-r")
	osVersion, err := cmd.Output()
	tools.PanicIfErr(err)

	winMgr := os.Getenv("XDG_CURRENT_DESKTOP")
	if len(winMgr) == 0 {
		fmt.Println("Warning! Cannot find window manager!")
	} else {
		winMgr += ","
	}
	
	osInfo, err := ioutil.ReadFile("/etc/os-release")
	tools.PanicIfErr(err)
	distro := ""

	for _, line := range strings.Split(string(osInfo), "\n") {
		if len(line) == 0 {
			continue
		}
		if line[:4] == "NAME" {
			distro = line[5:]
			distro = quotesReg.ReplaceAllString(distro, `\"`)
			break
		}
	}
			
	mask.SuperProperties = fmt.Sprintf(
		`{"os":"Linux","browser":"Discord Client","release_channel":"%v","client_version":"%v","os_version":"%v","os_arch":"x64","system_locale":"%v","window_manager":"%vunknown","distro":"%v","client_build_number":TODO:,"client_event_source":null}`, 
		releaseChan,
		mask.DiscordVersion,
		string(osVersion),
		mask.Locale,
		winMgr,
		distro,
	)

	mask.EncodeSuperProps()
}