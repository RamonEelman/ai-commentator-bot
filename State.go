package main

import "github.com/bwmarrin/discordgo"

type State struct {
	isConnected bool
	connection  discordgo.VoiceConnection
}
