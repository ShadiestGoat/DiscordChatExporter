//go:build windows

package config

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
	"golang.org/x/sys/windows/registry"
)

func (mask *HeadersMask) Auto() {
	cmdLoc := exec.Command("cmd", `Get-Culture | select -exp Name`)
	localeInp, err := cmdLoc.Output()
	tools.PanicIfErr(err)
	mask.Locale = string(localeInp)
	if len(mask.Locale) == 0 {
		fmt.Println("Warning! Could not auto pull the locale! Please report this as an issue. Will be using a default of en-US")
		mask.Locale = "en-US"
	}

	mask.PullDiscordVers(os.Getenv("APPDATA"))
	mask.UserAgent = fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/%v Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36", mask.DiscordVersion)

	releaseChan := "stable"

	if mask.UseCanary {
		releaseChan = "canary"
	}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	tools.PanicIfErr(err)
	defer k.Close()
	cmav, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	tools.PanicIfErr(err)
	cmiv, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	tools.PanicIfErr(err)
	cbn, _, err := k.GetStringValue("CurrentBuildNumber")
	tools.PanicIfErr(err)
	winVersion := fmt.Sprintf("%v.%v.%v", cmav, cmiv, cbn)

	mask.SuperProperties = fmt.Sprintf(
		`{"os":"Windows","browser":"Discord Client","release_channel":"%v","client_version":"%v","os_version":"%v","os_arch":"x64","system_locale":"%v","client_build_number"%v,"client_event_source":null}`,
		releaseChan,
		mask.DiscordVersion,
		winVersion,
		mask.Locale,
		mask.PullBuildId(),
	)
}
