package main

import (
	"log"
	"os"
)

type AiCommentatorConfig struct {
	guildId      string
	discordToken string
}

func initConfig() *AiCommentatorConfig {
	token, isAvailable := os.LookupEnv("DISCORD_TOKEN")
	if !isAvailable {
		log.Fatal("Could not find env variable discordToken")
	}
	guildId, isAvailable := os.LookupEnv("GUILD_ID")
	if !isAvailable {
		log.Fatal("Could not find env variable guildId")
	}
	return &AiCommentatorConfig{
		guildId:      guildId,
		discordToken: token,
	}
}
