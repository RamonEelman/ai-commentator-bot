package main

import (
	"github.com/Allan-Nava/fakeyou.go/configuration"
	"github.com/Allan-Nava/fakeyou.go/fakeyou"
	"log"
	"time"
)

var fakeYou = fakeyou.NewFakeYou(&configuration.Configuration{
	IsDebug: false,
	BaseUrl: "https://api.fakeyou.com/",
})

const trumpModel = "TM:djceg00wmcv5"

func voiceGeneratorTask(textChannel chan string) chan byte {
	voicesChannel := make(chan byte)
	go generateVoice(textChannel, voicesChannel)
	return voicesChannel
}

func generateVoice(textChannel chan string, voiceChannel chan byte) {
	for {
		text := <-textChannel
		voices, err := fakeYou.GenerateTTSAudio(text, trumpModel)
		if err != nil {
			log.Println("Could not retrieve generate a trump voice")
			time.Sleep(10 * time.Second)
		}
		pollJobUntilFinished(voices.InferenceJobToken)
	}
}

func pollJobUntilFinished(token string) {
	for {
		time.Sleep(5 * time.Second)
		response, err := fakeYou.PollTTSRequest(token)
		if err != nil {
			log.Println("Could not poll tts job", err)
			continue
		}
		if response.State.Status == "completed_status" {
			log.Println("Generating success: ")
			response2, err2 := fakeYou.GetListOfVoiceCategories()
			response.State
		}
	}
}
