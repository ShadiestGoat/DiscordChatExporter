package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"github.com/ShadiestGoat/DiscordChatExporter/config"
	"github.com/ShadiestGoat/DiscordChatExporter/tools"
)

type ConfigType config.Config


type NewDMChannelResp struct {
	Id string `json:"id"`
}

func (conf ConfigType) checkToken() {
	resp := []byte{}
	err := conf.discordFetch("/users/@me", &resp)
	if errors.Is(err, ErrBadAuth) {
		panic("Warning! This token is invalid!")
	} else {
		tools.PanicIfErr(err)
	}
	user := Author{}
	json.Unmarshal(resp, &user)
	fmt.Printf("Token is valid! Logged in as %v#%v!\n", user.Name, user.Discriminator)
}


func (conf ConfigType) FetchMain() {
	conf.checkToken()

	ParsedMessages := map[string]JSONMetaData{}
	NumLeft := conf.Filter.NumMax
	
	maxTime := conf.Filter.MaxTime
	minTime := conf.Filter.MinTime
	maxMsg := Message{}
	minMsg := Message{}

	switch conf.IdType {
		case config.ID_TYPE_CHANNEL:
		case config.ID_TYPE_GUILD:
			parsedGuilds := map[string]bool{}

			for _, id := range conf.Ids {
				if _, ok := parsedGuilds[id]; ok {
					continue
				}
				parsedGuilds[id] = true
			}

			conf.Ids = []string{}
			for _, guildId := range parsedGuilds {
				resp := []byte{}
				err := conf.discordFetch(fmt.Sprintf("/guilds/%v/channels", guildId), &resp)
				if errors.Is(err, Err404) {
					fmt.Printf("Warning! %v is not a guild id! Ignoring...", guildId)
					continue
				} else {
					tools.PanicIfErr(err)
				}
				allChannels := []Channel{}
				json.Unmarshal(resp, &allChannels)
				chanIds := []string{}
				for _, channel := range allChannels {
					switch channel.Type {
						case CHANNEL_TYPE_GUILD_CATEGORY, 
							 CHANNEL_TYPE_GUILD_STAGE_VOICE, 
							 CHANNEL_TYPE_GUILD_STORE, 
							 CHANNEL_TYPE_GUILD_VOICE:
							 	continue
					}
					if channel.Nsfw && conf.IgnoreNsfw {
						continue
					}
					chanIds = append(chanIds, channel.Id)
				}
				conf.Ids = append(conf.Ids, chanIds...)
			}
			
		case config.ID_TYPE_USER:
			newIds := []string{}
			for _, id := range conf.Ids {
				resBody := []byte{}
				err := conf.discordRequest(http.MethodPost, "/users/@me/channels", strings.NewReader(fmt.Sprintf(`{"recipient_id":"%v"}`, id)),&resBody)
				if errors.Is(err, Err404) {
					fmt.Printf("Warning! %v is either not a user, or a dm cannot be opened with them. You have to manually get the id from them!\n", id)
				} else {
					tools.PanicIfErr(err)
				}
				resp := NewDMChannelResp{}
				err = json.Unmarshal(resBody, &resp)
				tools.PanicIfErr(err)
				newIds = append(newIds, resp.Id)
			}
			conf.Ids = newIds
	}
	
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
			resp := []byte{}
			err := conf.discordFetch(fmt.Sprintf("/channels/%v/messages?limit=50%v", channel, beforeStr), &resp)
			tools.PanicIfErr(err)
			fmt.Println(string(resp))
			// parsed := []Message{}
			// json.Unmarshal(resp, &parsed)
			// messages = append(messages, parsed...)
			// done = true
		}
	}
	gotten, err := json.Marshal(ParsedMessages)
	tools.PanicIfErr(err)
	file, err := os.Create("gotten.json")
	tools.PanicIfErr(err)
	defer file.Close()
	file.Write(gotten)
}