package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ShadiestGoat/DiscordChatExporter/config"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

const BASE = "https://discordapp.com/api/v9"

type SingleMessageSearch struct {
	Msgs []Message `json:"messages"`
}

var ErrMsgNotFound = errors.New("msg not found")
var Err404 = errors.New("404")

func (conf ConfigType) discordFetch(uri string, body *[]byte) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v%v", BASE, uri), nil)
	tools.PanicIfErr(err)
	req.Header.Set("authorization", conf.Token)
	client := http.Client{
		Timeout: time.Second * 20,
	}
	res, err := client.Do(req)
	tools.PanicIfErr(err)
	panic(res.StatusCode)
	resBody, err := ioutil.ReadAll(res.Body)
	tools.PanicIfErr(err)
	*body = resBody
	return nil
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

type ConfigType config.Config

type JSONMetaData struct {
	Msgs []Message `json:"messages"`
	IDToIndex map[string]int `json:"idToIndex"`
	ByAuthor map[string][]string `json:"byAuthor"`
	Attachments []JSONAttachment `json:"attachments"`
	AuthorAttachment map[string]int `json:"attachment_byAuthor"`
}

type JSONAttachment struct {
	Attachment
	AuthorID string
}

func (conf ConfigType) FetchAll() {
	ParsedMessages := map[string]JSONMetaData{}
	NumLeft := conf.Filter.NumMax
	
	maxTime := conf.Filter.MaxTime
	minTime := conf.Filter.MinTime
	maxMsg := Message{}
	minMsg := Message{}
	
	if len(conf.Ids) == 1 {
		needChecking := 0
		if conf.Filter.MaxId != "" {
			maxMsgFetched, err := conf.FetchMsgId(conf.Ids[0], conf.Filter.MaxId)
			maxMsg = maxMsgFetched
			tools.PanicIfErr(err)
			needChecking++
			if conf.Filter.MinTime != 0 && conf.Filter.MinTime >= maxMsg.Timestamp {
				panic("BEFORE_ID is after BEFORE_TIME")
			}
		}
		if conf.Filter.MinId != "" {
			minMsgFetched, err := conf.FetchMsgId(conf.Ids[0], conf.Filter.MinId)
			tools.PanicIfErr(err)
			minMsg = minMsgFetched
			needChecking++
			if conf.Filter.MaxTime != 0 && minMsg.Timestamp >= conf.Filter.MaxTime {
				panic("AFTER_ID is after AFTER_TIME!")
			}
		}
		if needChecking == 2 {
			if minMsg.Timestamp > maxMsg.Timestamp {
				panic("Your BEFORE_ID is before AFTER_ID!")
			}
		}
	} else {
		if conf.Filter.MaxId != "" || conf.Filter.MinId != "" {
			fmt.Println("Alert! There is a Max ID or a Min ID, but multiple channels. This is not supported. If you want between certian dates on all channels, use BEFORE_TIME or AFER_TIME")
		}
	}

	if maxTime > maxMsg.Timestamp {
		maxTime = maxMsg.Timestamp
	}

	if minTime < minMsg.Timestamp {
		minTime = minMsg.Timestamp
	}

	for _, channel := range conf.Ids {
		if NumLeft == 0 {
			break
		}
		done := false
		messages := []Message{}
		for !done {
			if NumLeft == 0 {
				break
			}
			beforeStr := ""
			if len(messages) != 0 {
				beforeStr = fmt.Sprintf("&before=%v", messages[len(messages)-1].ID)
			}
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/channels/%v/messages?limit=100%v", BASE, channel, beforeStr), nil) // limit=2, because its non inclusive ://
			tools.PanicIfErr(err)
			req.Header.Set("authorization", conf.Token)
			client := http.Client{
				Timeout: time.Second * 20,
			}
			res, err := client.Do(req)
			tools.PanicIfErr(err)
			resBody, err := ioutil.ReadAll(res.Body)
			tools.PanicIfErr(err)
			parsed := []Message{}
			json.Unmarshal(resBody, &parsed)
			messages = append(messages, parsed...)
			done = true
		}
	}
	gotten, err := json.Marshal(ParsedMessages)
	tools.PanicIfErr(err)
	file, err := os.Create("gotten.json")
	tools.PanicIfErr(err)
	defer file.Close()
	file.Write(gotten)
}