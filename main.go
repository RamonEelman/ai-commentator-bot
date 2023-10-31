package main

import (
	"aicommentator/chatgpt"
	"aicommentator/mongo"
	"fmt"
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
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	defer dg.Close()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	matchSummaryChannel := mongo.StartChangeListener()
	commentateResultChannel := chatgpt.StartChatGPTTask(matchSummaryChannel)
	voiceChannel := RunVoiceTTSTask(commentateResultChannel)
	eventChannel := RunVoiceEventHandler(dg, config)
	RunVoiceConnectTask(dg, config, eventChannel)
	RunVoiceSendTask(voiceChannel, eventChannel)
}

func RunVoiceSendTask(channel chan chan []byte, eventChannel chan VoiceEvent) {
	for {
		newVoiceMessage := <-channel
		eventChannel <- VoiceEvent{
			eventType:    PlaySound,
			voiceMessage: newVoiceMessage,
		}
	}
}
