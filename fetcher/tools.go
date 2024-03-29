package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/discord"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

const DOMAIN = "discord.com"
const API_BASE = "/api/v9"

var ErrMsgNotFound = errors.New("msg not found")
var Err404 = errors.New("404")
var ErrBadAuth = errors.New("error 401: unauthorized! this means the token is bad")
var jar *cookiejar.Jar

type RateLimit struct {
	RetryAfter float64 `json:"retry_after"`
}

func DownloadMedia(mediaDir string, urlToFetch string, name string) {
	resp, err := http.Get(urlToFetch)

	if err != nil {
		if errors.Unwrap(err).Error()[:8] != "dial tcp" {
			panic(err)
		} else {
			return
		}
	}

	file, err := os.Create(filepath.Join(mediaDir, name)) // Will overwrite it no matter what, so embed image duplicates are fine :)
	tools.PanicIfErr(err)

	dwMedia, err := io.ReadAll(resp.Body)
	tools.PanicIfErr(err)

	defer file.Close()

	file.Write(dwMedia)
}

func (conf ConfigType) GetCookie() {
	req, err := http.NewRequest("GET", "https://"+conf.HeadersMask.DomainPrefix+"discord.com", nil)
	tools.PanicIfErr(err)
	jarN, err := cookiejar.New(nil)
	tools.PanicIfErr(err)
	jar = jarN

	req.Header.Set("User-Agent", conf.HeadersMask.UserAgent)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}

	_, err = client.Do(req)
	tools.PanicIfErr(err)
}

func (conf ConfigType) discordRequest(method string, uri string, body io.Reader, respBody *[]byte) error {
	urlBase := "https://" + conf.HeadersMask.DomainPrefix + DOMAIN
	req, err := http.NewRequest(method, urlBase+API_BASE+uri, body)
	tools.PanicIfErr(err)
	req.Header.Set("Authorization", conf.Token)
	req.Header.Set("User-Agent", conf.HeadersMask.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authority", conf.HeadersMask.DomainPrefix+DOMAIN)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("X-Discord-Locale", conf.HeadersMask.Locale)
	req.Header.Set("X-Super-Properties", conf.HeadersMask.SuperProperties)
	req.Header.Set("Accept-Language", conf.HeadersMask.Locale)
	// referer... I don't think you need it since its a more standard header..? Idk For the record its like this: https://canary.discord.com/channels/@me/CHANID

	client := http.Client{
		Timeout: time.Second * 20,
		Jar:     jar,
	}

	res, err := client.Do(req)
	tools.PanicIfErr(err)
	resBody, err := io.ReadAll(res.Body)
	tools.PanicIfErr(err)

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return Err404
		} else if res.StatusCode == 401 {
			return ErrBadAuth
		} else if res.StatusCode == 429 {
			rates := RateLimit{}
			json.Unmarshal(resBody, &rates)
			fmt.Printf("Rate limit achieved! Retrying in %vs...\n", rates.RetryAfter)
			time.Sleep(time.Duration((rates.RetryAfter + 0.1) * float64(time.Second))) // 0.1 just in case
			err = conf.discordRequest(method, uri, body, respBody)
			tools.PanicIfErr(err)
			return nil
		} else {
			panic(fmt.Sprintf("Unknown status code %v detected!", res.StatusCode))
		}
	}
	*respBody = resBody

	return nil
}

func (conf ConfigType) discordFetch(uri string, respBody *[]byte) error {
	return conf.discordRequest(http.MethodGet, uri, nil, respBody)
}

func (conf ConfigType) FetchMsgId(channel string, id string) (discord.Message, error) {
	resBody := []byte{}
	err := conf.discordFetch(fmt.Sprintf("/channels/%v", channel), &resBody) // we don't actuall care about output so it's fine to use resBody
	if errors.Is(err, Err404) {
		panic(fmt.Sprintf("CHANNEL '%v' was not found", channel))
	} else {
		tools.PanicIfErr(err)
	}

	conf.discordFetch(fmt.Sprintf("/channels/%v/messages?around=%v&limit=1", channel, id), &resBody)

	if string(resBody) == "[]" {
		return discord.Message{}, ErrMsgNotFound
	}

	msgs := []discord.Message{}
	err = json.Unmarshal(resBody, &msgs)
	tools.PanicIfErr(err)
	return msgs[0], nil
}
