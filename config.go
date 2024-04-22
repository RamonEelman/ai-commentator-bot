package main

import (
	"log"
	"os"
)

type AiCommentatorConfig struct {
	guildId      string
	discordToken string
	chatChannel  string
}

func initConfig() *AiCommentatorConfig {
	token, isAvailable := os.LookupEnv("DISCORD_TOKEN")
	if !isAvailable {
		log.Fatal("Could not find env variable DISCORD_TOKEN")
	}
	guildId, isAvailable := os.LookupEnv("GUILD_ID")
	if !isAvailable {
		log.Fatal("Could not find env variable GUILD_ID")
	}
	chatChannel, isAvailable := os.LookupEnv("MESSAGE_CHANNEL")
	if !isAvailable {
		log.Fatal("Could not find env variable MESSAGE_CHANNEL")
	}

	return &AiCommentatorConfig{
		guildId:      guildId,
		discordToken: token,
		chatChannel:  chatChannel,
	}
}
