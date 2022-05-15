package config

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
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
	} else if !dExist && !dCExist {
		panic("Could not detect discord!")
	} else if dExist && !dCExist && mask.UseCanary && len(os.Getenv("USE_CANARY")) != 0 {
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
			break
		}
	}
	if len(mask.DiscordVersion) == 0 {
		panic("Couldn't auto-pull the discord version!")
	}
}

var discordAssetReg = regexp.MustCompile(`assets/+?([a-z0-9]+?)\.js`)
var buildNumReg = regexp.MustCompile(`buildNumber`)

// https://github.com/Merubokkusu/Discord-S.C.U.M/blob/master/discum/start/superproperties.py thank you!
func (mask HeadersMask) PullBuildId() string {
	respRaw, err := http.Get(fmt.Sprintf("https://%vdiscord.com/login", mask.DomainPrefix)) //idk if its different but i am very tired atm so yk
	tools.PanicIfErr(err)
	resp, err := ioutil.ReadAll(respRaw.Body)
	tools.PanicIfErr(err)
	assets := discordAssetReg.FindAll(resp, 50)
	buildFileRaw, err := http.Get(fmt.Sprintf("https://%vdiscord.com/%v", mask.DomainPrefix, string(assets[len(assets)-2])))
	tools.PanicIfErr(err)
	buildFile, err := ioutil.ReadAll(buildFileRaw.Body)
	tools.PanicIfErr(err)

	buildFileS := string(buildFile)

	buildLoc := buildNumReg.FindStringIndex(buildFileS)

	if len(buildLoc) == 0 {
		panic("Build num not located! Auto loading has failed :(")
	}

	return buildFileS[buildLoc[1]+3 : buildLoc[1]+9]
}

func (mask *HeadersMask) EncodeSuperProps() {
	mask.SuperProperties = base64.StdEncoding.EncodeToString([]byte(mask.SuperProperties))
}

var localeReg = regexp.MustCompile(`.._..`)
var versionReg = regexp.MustCompile(`\d+\.\d+\.\d+`)
var quotesReg = regexp.MustCompile(`"`)
