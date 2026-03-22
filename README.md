# What's the Score

## Usage

```
What's the Score: Check the scores without leaving your terminal.

List today's games around the league:
  wts ls

Check on a specific team:
  wts ls lakers

WTS requires an API key to BALLDONTLIE. Sign up for a free key at
https://balldontlie.io/ and save it using:
  wts set-api-key

Usage:
  wts [command]

Available Commands:
  list        List scores by league or team
  set-api-key Set your BALLDONTLIE API key

Flags:
      --config string   config file (default is $HOME/.config/wts.toml)
      --debug           Enable debug logging
  -h, --help            help for wts

Use "wts [command] --help" for more information about a command.
```

## Installation

### Homebrew

```
brew tap andyyoon/whats-the-score
brew install whats-the-score
wts
```

### Download

Download the latest binary for your OS from the releases page. https://github.com/andyyoon2/whats-the-score/releases

