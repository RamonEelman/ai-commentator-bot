package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type VoiceEvent struct {
	eventType    VoiceEventType
	channelId    string
	voiceMessage chan []byte
}

type VoiceEventType = int8

const (
	Connect    VoiceEventType = 0
	Disconnect VoiceEventType = 1
	PlaySound  VoiceEventType = 2
)

func RunVoiceEventHandler(session *discordgo.Session, config *AiCommentatorConfig) chan VoiceEvent {
	state := &State{
		isConnected: false,
	}
	events := make(chan VoiceEvent, 4)
	go func() {
		for {
			event := <-events
			switch {
			case event.eventType == Connect:
				connectToVoice(session, event.channelId, config.guildId, state)
			case event.eventType == Disconnect:
				disconnectFromVoice(state)
			case event.eventType == PlaySound:
				if state.isConnected {
					playSound(state, event.voiceMessage)
				} else {
					drainToVoid(event.voiceMessage)
				}
			}
		}
	}()
	return events
}

func drainToVoid(channel chan []byte) {
	for _ = range channel {
	}
}

func playSound(session *State, channel chan []byte) {
	drainChannelBlocking(channel, session.connection.OpusSend)
}

func drainChannelBlocking(src <-chan []byte, dst chan<- []byte) {
	// Read from the source channel and write into the destination channel
	for chunk := range src {
		dst <- chunk
	}
}

func disconnectFromVoice(state *State) {
	if state.isConnected == false {
		return
	}
	err := state.connection.Disconnect()
	if err != nil {
		log.Fatal("Error disconnecting from voice chat")
	}
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
