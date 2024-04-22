package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func RunChatSendTask(dg *discordgo.Session, config *AiCommentatorConfig) chan<- string {
	chatChannel := make(chan string)
	go retrieveTextAndSendDiscordTask(chatChannel, dg, config.chatChannel)
	return chatChannel
}

func retrieveTextAndSendDiscordTask(channel chan string, dg *discordgo.Session, chatChannelId string) {
	for {
		message := <-channel
		_, err := dg.ChannelMessageSend(chatChannelId, message)
		if err != nil {
			log.Println("Error sending text message: '", message, "' with error: ", err)
		}
	}

}
