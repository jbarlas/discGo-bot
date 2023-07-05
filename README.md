# DiscGo Bot

Simple Discord bot to learn go!

## Set Up
### Discord Bot
1. Create an application using the [Discord Developer Portal](https://discord.com/developers)
2. Create a bot and log its token (to add in `config.json`)
3. Add the bot to your server via OAuth2 -> URL Generator -> bot
  - Give your bot the proper permissions (all text permissions + read messages)
  - Be sure to authorize `MESSAGE CONTENT INTENT` on the bot home page
### Config
Ensure there is a config.json file in the parent directory which has the following form:
```json
{
    "Token": "<botToken>",
    "BotPrefix": "!"
}
```

## To Run
`go run main.go`