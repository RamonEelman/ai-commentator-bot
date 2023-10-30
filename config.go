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
	token, isAvailable := os.LookupEnv("discordToken")
	if !isAvailable {
		log.Fatal("Could not find env variable discordToken")
	}
	guildId, isAvailable := os.LookupEnv("guildId")
	if !isAvailable {
		log.Fatal("Could not find env variable guildId")
	}
	return &AiCommentatorConfig{
		guildId:      guildId,
		discordToken: token,
	}
}
