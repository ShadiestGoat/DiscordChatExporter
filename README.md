# Discord Chat exporter

<p align="center">
    <img src="logo-m.svg" width="25%">
</p>

> :warning: This project is out only as a beta, if you see any bugs, please make an issue :warning:
> :warning: Due to this being in beta, API changes are to be expected, even if semver would disagree with this change :warning:

This is a golang implemintation of a discord chat explorer like [this](https://github.com/Tyrrrz/DiscordChatExporter) or [this](https://github.com/mahtoid/DiscordChatExporterPy).

## Features

- Highly configurable 
- Barely any dependencies
- Multiple ways of exporting (TEXT (ie. a log), JSON w/ meta info, HTML (visual))
- Allows auto parsing of channels
- Allows filters that aren't mutually exclusive (unlike discord)
- Golang superiority
- Mass downloads from guild(s)
- Commitments to less bloat
- Customisability 

## Setup

Please check the [wiki](https://github.com/ShadiestGoat/DiscordChatExporter/wiki/Setup)

## Faq

### Why a user token?

You have to use a user token, since bots are no longer allowed to download messages later than the last 2 weeks. 

### Discord TOS?

This is technically breaking the discord terms of service, as it is user botting. 

However, so far, haven't been banned so make your own conclusions.

## Contributing

Pull requests are always welcome! Just make sure to follow the code of conduct :)

## Roadmap


- [ ] Make this a more cli friendly app
    - [ ] Add an option for command line flags for config
    - [ ] Add a global location for themes (probs would have to embed the default ones)
- [ ] Autoparsing of IDs
- [ ] Writing tests
- [ ] html:
    - [ ] Add system messages (like calls or pins)
    - [ ] Improving html
    - [ ] Adding more themes
    - [ ] Add some reactivity (eg. actially go to the message being replied to when clicked)
- [ ] Add a GUI
