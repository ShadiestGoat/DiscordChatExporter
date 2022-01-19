package config

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

func (mask *HeadersMask) PullDiscordVers(discordPath string) {
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

	if dExist && !dCExist && mask.UseCanary && len(os.Getenv("USE_CANARY")) != 0 {
		fmt.Println("Warning! Stable discord found but canary not found, and the preferance is for canary! We will be using false for USE_CANARY for this download")
		mask.UseCanary = false
	} else if !mask.UseCanary && dCExist && !dExist && len(os.Getenv("USE_CANARY")) != 0 {
		fmt.Println("Warning! Canary discord found but stable not found, and the preferance is for stable! We will be using true for USE_CANARY for this download")
		mask.UseCanary = true
	}
	
	if mask.UseCanary {
		mask.DomainPrefix = "canary."
	}

	discordFiles := normalD

	if mask.UseCanary {
		discordFiles = normalCD
	}

	for _, fileInfo := range discordFiles {
		if !fileInfo.IsDir() {
			continue
		}
		name := fileInfo.Name()
		if versionReg.MatchString(name) {
			mask.DiscordVersion = name
		}
	}
}

func (mask *HeadersMask) EncodeSuperProps() {
	mask.SuperProperties = base64.StdEncoding.EncodeToString([]byte(mask.SuperProperties))
}

var localeReg = regexp.MustCompile(`.._..`)
var versionReg = regexp.MustCompile(`\d+\.\d+\.\d+`)
var quotesReg = regexp.MustCompile(`"`)