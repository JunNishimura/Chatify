<h1 align='center'>
  ♫ Chatify ♫ <br/>chat-based music recommendation tool
</h1>

<p align='center'>
  <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/JunNishimura/Chatify">
  <img alt="GitHub" src="https://img.shields.io/github/license/JunNishimura/Chatify">
  <a href="https://goreportcard.com/report/github.com/JunNishimura/Chatify"><img src="https://goreportcard.com/badge/github.com/JunNishimura/Chatify" alt="Go Report Card"></a>
</p>

![screen](https://github.com/JunNishimura/Chatify/assets/28744711/396d90a8-eeb6-4f04-a705-94f9a1c6e733)


# 📖 Overview
Chatify is a TUI(Terminal User Interface) tool that combines the OpenAI API with the Spotify API, allowing an AI bot to recommend the music you are looking for through conversation. 

1. You can find the music you are looking for by answering questions.
2. You can listen to the music by selecting the track recommended by chatify.
3. You can store results of recommendations as playlists.

# ⚠️ Notice
Unfortunately, Chatify is currently not free to use; you need to have a Spotify Premium account and pay to use the OpenAI API.

# 👜 Prerequisites
Chatify requires two things. 
1. You need to have [a Spotify Premium Account](https://www.spotify.com/premium/) to use Spotify API.
2. You need to have [a OpenAI account](https://platform.openai.com/login) to use OpenAI API.

# 💻 How to use
**It is recommended that the terminal is set to full-screen size when to use Chatify.**

## 1. Install
### Homebrew Tap
```
brew install JunNishimura/tap/Chatify
```
### Go install
```
go install github.com/JunNishimura/Chatify@latest
```

## 2. Preparation
Before saying hi to Chatify, please prepare the following three items.

```
1. Spotify App Client ID
2. Spotify App Client Secret
3. OpenAI API key
```

### Spotify API

Please create any app from the [Spotify for Developers Dashboard](https://developer.spotify.com/dashboard) and retrieve the Client ID and Secret from the settings screen.

### OpenAI API

Please create an API key from [the OpenAI account screen](https://platform.openai.com/account/api-keys) and obtain it.

## 3. Greetings to Chatify
You need to provide a couple of information to Chatify at first.

```
$ chatify greeting
```

## 4. Talk with Chatify
Let's talk to Chatify and embark on a journey to discover new music! 

```
$ chatify hey
```

# ⌨️ Operation
| Key | Action |
| ---- | ---- |
| tab | switch the view |
| characters | input words |
| q, ctrl+c | quit |
| enter | answer, select |
| ↑, ↓ | move in list |
| ←, → | turn page |

# 🔨 Options
## `greeting` command
###  `-p`, `--port`
Flag to specify the port number for Spotify authorization. Default number is 8888.

## `hey` command
### `-n`, `--number`
Flag to specify the number of recommendations. Default number is 25 and maximum number is 100.

### `-p`, `--playlist`
Flag to enable chatify to make playlist based on the recommendation. Default is false.


# 🪧 License
Chatify is released under MIT License. See [MIT](https://raw.githubusercontent.com/JunNishimura/Chatify/main/LICENSE)
