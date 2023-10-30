package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type VoiceEvent struct {
	eventType VoiceEventType
	channelId string
}

type VoiceEventType = int8

const (
	Connect    VoiceEventType = 0
	Disconnect VoiceEventType = 1
	PlaySound  VoiceEventType = 2
)

func RunVoiceEventHandler(events chan VoiceEvent, session *discordgo.Session, config *AiCommentatorConfig) {
	state := &State{
		isConnected: false,
	}
	go func() {
		for {
			event := <-events
			switch {
			case event.eventType == Connect:
				connectToVoice(session, event.channelId, config.guildId, state)
			case event.eventType == Disconnect:
				disconnectFromVoice(state)
			case event.eventType == PlaySound:
				playSound()
			}
		}
	}()
}

func playSound() {
	// to Do play sound
}

func disconnectFromVoice(state *State) {
	if state.isConnected == false {
		return
	}
	state.connection.Close()
	state.isConnected = false
}

func connectToVoice(session *discordgo.Session, cid string, gid string, state *State) {
	if state.isConnected == true {
		return
	}
	join, err := session.ChannelVoiceJoin(gid, cid, false, false)
	if err != nil {
		log.Println("could not connect to voice", err)
		return
	}
	state.isConnected = true
	state.connection = *join
}
