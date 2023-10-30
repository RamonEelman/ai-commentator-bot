package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"time"
)

func RunVoiceConnectTask(session *discordgo.Session, config *AiCommentatorConfig, voiceEvents chan VoiceEvent) {
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case _ = <-ticker.C:
				guildEvents, err := session.GuildScheduledEvents(config.guildId, false)
				if err != nil {
					log.Println("Error could not find guild guildEvents", err)
				}
				isActive, event := isEventActive(guildEvents)
				switch isActive {
				case true:
					voiceEvents <- VoiceEvent{
						eventType: Connect,
						channelId: event.ChannelID,
					}
				case false:
					voiceEvents <- VoiceEvent{
						eventType: Disconnect,
					}
				}
			}
		}
	}()
}

func isEventActive(events []*discordgo.GuildScheduledEvent) (bool, *discordgo.GuildScheduledEvent) {
	for _, event := range events {
		pavlovRegex := regexp.MustCompile("(?i)pavlov")
		if event.Status == discordgo.GuildScheduledEventStatusActive &&
			(pavlovRegex.MatchString(event.Name) || pavlovRegex.MatchString(event.Description)) {
			return true, event
		}
	}
	return false, nil
}
