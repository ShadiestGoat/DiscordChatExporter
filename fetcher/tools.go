package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

func (author Author) AvatarUrl() string {
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/376079696489742338/%v.webp?size=512", author.Avatar)
}

const BASE = "https://discordapp.com/api/v9"

type SingleMessageSearch struct {
	Msgs []Message `json:"messages"`
}

var ErrMsgNotFound = errors.New("msg not found")
var Err404 = errors.New("404")

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
	if res.StatusCode == 404 {
		return Err404
	}
	resBody, err := ioutil.ReadAll(res.Body)
	tools.PanicIfErr(err)
	*respBody = resBody
	return nil
}

func (conf ConfigType) discordFetch(uri string, respBody *[]byte) error {
	return conf.discordRequest(http.MethodGet, uri, nil, respBody)
}


func (conf ConfigType) FetchMsgId(channel string, id string) (Message, error) {
	resBody := []byte{}
	err := conf.discordFetch(fmt.Sprintf("/channels/%v", channel), &resBody) // we don't actuall care about output so it's fine to use resBody
	if errors.Is(err, Err404) {
		panic(fmt.Sprintf("CHANNEL '%v' was not found", channel))
	} else {
		tools.PanicIfErr(err)
	}

	err = conf.discordFetch(fmt.Sprintf("/channels/%v/messages?around=%v&limit=1", channel, id), &resBody)
	if string(resBody) == "[]" {
		return Message{}, ErrMsgNotFound
	}
	msgs := SingleMessageSearch{}
	err = json.Unmarshal(resBody, &msgs)
	tools.PanicIfErr(err)
	return msgs.Msgs[0], nil
}