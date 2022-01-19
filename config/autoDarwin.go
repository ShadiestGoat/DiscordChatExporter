//go:build darwin

package config

import (
	"fmt"
	"os"
	"os/exec"

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
	mask.PullDiscordVers("~/Library/Application Support")
	mask.UserAgent = fmt.Sprintf("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) discord/%v Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36", mask.DiscordVersion)
	
	releaseChan := "stable"

	if mask.UseCanary {
		releaseChan = "canary"
	}
	
	cmd := exec.Command("uname", "-r")
	osVersion, err := cmd.Output()
	tools.PanicIfErr(err)

	mask.SuperProperties = fmt.Sprintf(
		`{"os":"Mac OS X","browser":"Discord Client","release_channel":"%v","client_version":"%v","os_version":"%v","os_arch":"x64","system_locale":"%v","client_build_number":TODO:,"client_event_source":null}`,
		releaseChan,
		mask.DiscordVersion,
		string(osVersion),
		mask.Locale,
	)
	mask.EncodeSuperProps()
}