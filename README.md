<h1 align='center'>
  ♫ Chatify ♫ <br/>Chat-based music recommendation tool
</h1>

<p align='center'>
  <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/JunNishimura/Chatify">
  <img alt="GitHub" src="https://img.shields.io/github/license/JunNishimura/Chatify">
  <a href="https://goreportcard.com/report/github.com/JunNishimura/Chatify"><img src="https://goreportcard.com/badge/github.com/JunNishimura/Chatify" alt="Go Report Card"></a>
</p>

## 📖 Overview
Chatify is a TUI(Terminal User Interface) tool that combines the OpenAI API with the Spotify API, allowing an AI bot to recommend the music you are looking for through conversation. 

## ⚠️ Notice
Unfortunately, Chatify is currently not free to use; you need to have a Spotify Premium account and pay to use the OpenAI API.

## 👜 Prerequisites
Chatify requires two things. 
1. You need to have [a Spotify Premium Account](https://www.spotify.com/premium/) to use Spotify API.
2. You need to have [a OpenAI account](https://platform.openai.com/login) to use OpenAI API.

## 💻 How to use
### 1. Install
#### Homebrew Tap
```
brew install JunNishimura/tap/Chatify
```
#### Go install
```
go install github.com/JunNishimura/Chatify@latest
```

### 2. Preparation
Before saying hi to Chatify, please prepare the following three items.

1. Spotify App Client ID
2. Spotify App Client Secret
3. OpenAI API key

As for `Spotify App Client ID` and `Spotify App Client Secret`, please create any app from the [Spotify for Developers Dashboard](https://developer.spotify.com/dashboard) and retrieve them from the settings screen.

For `OpenAI API key`, please create an API key from [the OpenAI account screen](https://platform.openai.com/account/api-keys) and obtain it.

### 3. Greetings to Chatify
You need to provide a couple of information to Chatify at first.

```
$ chatify greeting
```

### 3. Talk with Chatify
Let's talk to Chatify and embark on a journey to discover new music! 

```
$ chatify hey
```

## 🪧 License
Chatify is released under MIT License. See [MIT](https://raw.githubusercontent.com/JunNishimura/Chatify/main/LICENSE)
