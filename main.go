package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create a new Discord session using the provided bot token.
	config := initConfig()

	dg, err := discordgo.New("Bot " + config.discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.Identify.Intents = INTENT_CONFIG
	// Register the messageCreate func as a callback for MessageCreate events.
	voiceChannel := make(chan VoiceEvent)
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	RunVoiceConnectTask(dg, config, voiceChannel)
	RunVoiceEventHandler(voiceChannel, dg, config)
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	close(voiceChannel)
	// Cleanly close down the Discord session.
	dg.Close()
}
