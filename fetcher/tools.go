package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/discord"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

const BASE = "https://discordapp.com/api/v9"

var ErrMsgNotFound = errors.New("msg not found")
var Err404 = errors.New("404")
var ErrBadAuth = errors.New("error 401: unauthorized! this means the token is bad")

type RateLimit struct {
	RetryAfter float64 `json:"retry_after"`
}

func (conf ConfigType) discordRequest(method string, uri string, body io.Reader, respBody *[]byte) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", BASE, uri), body)
	tools.PanicIfErr(err)
	req.Header.Set("authorization", conf.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 20,
	}
	res, err := client.Do(req)
	tools.PanicIfErr(err)
	resBody, err := ioutil.ReadAll(res.Body)
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
			err = conf.discordRequest(method, uri, body, &resBody)
			if err != nil {
				return err
			}
		} else {
			panic(fmt.Sprintf("Unknown status code %v detected!", res.StatusCode))
		}
	}
	tools.PanicIfErr(err)
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