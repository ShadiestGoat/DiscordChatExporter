# Discord Chat exporter

> :warning: THIS PROJECT IS STILL UNDER DEVELOPMENT, NOT EVEN OUT AS A BETA :warning:

This is a golang implemintation of a discord chat explorer i found a while back & can't find again

## Features

- Highly configurable 
- No need for dependencies
- Multiple ways of exporting
- Allows DM channels to be automatically pulled (without id)
- Allows filters that aren't mutually exclusive
- Golang superiority
- Mass downloads from guild(s)
- Bloatless 

## Config

This is configurated through env variables. You can do this either from command line or from the .env file:

1. `TOKEN="..." ID="..." ./`
2. (in `.env`):

```js
const test = ""
TOKEN="..."
ID="..."
```

Here are all the available config options:

| Value | Accepted Values | Description | Default |
|:----:|:----:|:----:|:----:|
| IGNORE_ENV_FILE | Boolean Type | If true, it will not load the `.env` file. Does not work if it is in the `.env` file | false |
| ENV_FILENAME | Any string | The location of the .env file | ".env" |
| TOKEN | User Token Type | The **user** token that will be used to download the messages. For more info check faq! | :x: |
| ID_TYPE | "USER", "CHANNEL", "GUILD" | The type of the ID provided in the next line. If it is USER, it will be taken as a user id, and will pull DM conversations with them. If It's channel, it will download all the messages from that channel id. If type is GUILD, it will get all channels in the guild & export them. Most limits won't work if it's guild (except number of msgs) | "CHANNEL" |
| ID | Id Type | The Id of the channel or user (specified in above line). Please note that there can be multiple ids here, separated by a space. They all have to be the same type though (think all USER, CHANNEL, or GUILD) | :x: |
| DOWNLOAD_MEDIA | Boolean type | If true, then it downloads all the media & attachments locally, otherwise it just links to it | false |
| IGNORE_NSFW_CHANNEL | Boolean type | If true, then it will ignore all automatically fetched channels (ie. from guilds) | false |
| USE_LIMIT_50 | Boolean Type | If true, it uses limit=50 rather than 100 for message downloading. This is the discord default, so it is technically safer for your account. | false |
| EXPORT_TYPE | "TEXT", "JSON", "HTML" | The type to be exported as. For more info, check faq! ==TODO== | "JSON" |
| EXPORT_LOCATION | Any string | The location where this will be exported to. This will create a folder under that name that contains the media (see `DOWNLOAD_MEDIA`), and the messages. This is templated, and you can use variables. Variables are wrapped in `{{$}}` (eg. `{{$CHANNEL_ID}}`), vars are `CHANNEL_ID`. | "output/$CHANNEL_ID" (cross-platform) |
| EXPORT_HTML_THEME | Any string | The name of the theme to be used for discord. Pre made options are "light", "dark", "black". There is no guide on custom themes (yet), but look at premade themes | "dark" |
| MSG_LIMIT_NUM | any number, or "all" | The number of messages to download. If it's "all", then it'll download all the messages | "all" |
| EXPORT_JSON_TOOLS | Boolean Type | If true, it will export special tools in the json along with messages. `{"messages": [Array of messages], "idToIndex": {"MessageId": Index (Number)}, byAuthor: {"AuthorId": [Array of Msg Ids (string)]}, "attachments": [], "attachment_byAuthor": {"AuthorId": IndexOfAttachment (Number)}}`. If false, this will only output an array of messages. This is only taken into account if `EXPORT_TYPE` = "JSON" | true |
| EXPORT_PLAIN_FORMAT | Any string | A template for how each msg is presented. vars are wrapped in `{{$VAR_HERE}}` (eg. `{{$AUTHOR_ID}}`). Available variables are: `AUTHOR_NAME`, `AUTHOR_ID`, `TIMESTAMP` (utc timestamp of the msg), `IS_REPLY` ("true" if it is replying to a msg), `WAS_EDITED`, `CONTENT` (note, in has it has an attachment or sticker CONTENT will be the id/url), `HAS_ATTACHMENT`, `ATTACHMENT_URL` (a "," seperated list of urls, which are also enclosed in `""`, ie. "testtest","testtest2", but w/ urls), `IS_STICKER`, `STICKER_IDS` (a "," separated list of sticker IDs). A newline will be added regardless of anything at the end of each entry | string |
| BEFORE_ID | nil, Id type | Only export messages before this message id. Not mutually exclusive w/ AFTER_ID, however if this one is before AFTER_ID, it will throw errors | nil |
| AFTER_ID | nil, Id Type | Only export messages after this message id. Not mutually exclusive w/ BEFORE_ID, however if this one is after BEFORE_ID, it will throw errors | nil |
| BEFORE_TIME | nil, Timestamp | Only export messages before the given time period. Not mutually exclusive with anything, but could throw errors if it is after the AFTER_ID and AFTER_TIME | nil |
| AFTER_TIME | nil, Timestamp | Only export messages before the given time period. Not mutually exclusive with anything, but could throw errors if it is after the BEFORE_ID and BEFORE_TIME | nil |

Hera are all the types mentioned:

| Type | Description | Values |
| ---- | ---- | ---- |
| Boolean | Yes or no | Case insensitive, `true`: "true", "t", "yes", "y", "1". `false`: "false", "f", "yes", "y", "0" |
| User Token | A user token | Any string, that matches a user token reg (ie. starts w/ `mfa.`) |
| nil | An empty config option. Has no value, or empty string | Empty |
| Id | A discord Id (snowflake?) | A string of 18 digits |
| Timestamp | A Epoch Unix Timestamp | Epoch Unix Timestamp |

## Setup guide

There are several ways of setting this up. You can either build from source or use pre-built binaries found on the releases section. 

### Pre-Built

Download the binary from the releases section, with the correct operating system. 

### Build from source

1. Clone this repo (`git clone https://github.com/ShadiestGoat/DiscordChatExporter`)
2. Move into the directory `cd DiscordChatExporter`
3. Install dependencies `go get github.com/ShadiestGoat/DiscordChatExporter`
4. Build
	- Either build only for your operating system, `go build` (optionally use `GOOS=` + `windows`, `linux` or `darwin`)
	- Or build for all oses, `sh ./build.sh`

## Usage Guide

1. You need a user token ([Why?](#Why-a-user-token?)). [How to get one](https://www.just-fucking-google.it/?s=how%20to%20get%20discord%20user%20token&e=fingerxyz)
2. Set up your env (either `.env` or through command line) [details](#Config)
3. Open terminal & navigate to the binary you got (use `cd`)
4. Execute the binary. 


## Faq

### Why a user token?

You have to use a user token, since bots are no longer allowed to download messages later than the last 2 weeks. 

### Discord TOS?

This is technically breaking the discord terms of service, as it is user botting. 

However, so far, haven't been banned so make your own conclusions.

