package main

import (
	"github.com/Allan-Nava/fakeyou.go/configuration"
	"github.com/Allan-Nava/fakeyou.go/fakeyou"
	"github.com/jonas747/dca"
	"io"
	"log"
	"net/http"
	"time"
)

var fakeYou = fakeyou.NewFakeYou(&configuration.Configuration{
	IsDebug: false,
	BaseUrl: "https://api.fakeyou.com/",
})

const trumpModel = "TM:djceg00wmcv5"

func RunVoiceTTSTask(textChannel <-chan string) chan chan []byte {
	voicesChannel := make(chan chan []byte)
	go generateVoice(textChannel, voicesChannel)
	return voicesChannel
}

var settings = &dca.EncodeOptions{
	Volume:           "1.0",
	Channels:         1,
	FrameRate:        48000,
	FrameDuration:    20,
	Bitrate:          128,
	Application:      dca.AudioApplicationVoip,
	CompressionLevel: 10,
	PacketLoss:       1,
	BufferedFrames:   100, // At 20ms frames that's 2s
	VBR:              true,
	StartTime:        0,
}

func generateVoice(textChannel <-chan string, voiceChannel chan<- chan []byte) {
	for {
		text := <-textChannel
		println(text)
		voices, err := fakeYou.GenerateTTSAudio(text, trumpModel)
		log.Println(voices.InferenceJobToken)
		if err != nil {
			log.Println("Could not retrieve generate a trump voice", err)
			time.Sleep(10 * time.Second)
			continue
		}
		downloadPath := pollJobUntilFinished(voices.InferenceJobToken)
		if downloadPath != "" {
			reader := downloadAudioFile(downloadPath)
			encoderSession, err := dca.EncodeMem(reader, dca.StdEncodeOptions)
			if err != nil {
				log.Println("Could not create audio session from reader", err)
			}
			currentTTSChannel := make(chan []byte, 100)
			voiceChannel <- currentTTSChannel
			encoderToOpusCopy(encoderSession, currentTTSChannel)
			reader.Close()
		}
	}
}

func encoderToOpusCopy(encoderSession *dca.EncodeSession, currentTTSChannel chan<- []byte) {
	defer encoderSession.Cleanup()
	defer close(currentTTSChannel)
	for {
		frame, err := encoderSession.OpusFrame()
		switch err {
		case nil:
			currentTTSChannel <- frame
		case io.EOF:
			return
		default:
			log.Println("Encountered error reading audio frames", err)
			return
		}
	}
}

func pollJobUntilFinished(token string) string {
	for {
		time.Sleep(5 * time.Second)
		response, err := fakeYou.PollTTSRequest(token)
		if err != nil {
			log.Println("Could not poll tts job", err)
			continue
		}
		switch response.State.Status {
		case fakeyou.StateStatusPending, fakeyou.StateStatusStarted:
			log.Println("Generating with status:" + response.State.Status)

		case fakeyou.StateStatusCompleteSuccess:
			log.Println("Successfully created an audio")
			audioDownloadPath := "https://storage.googleapis.com/vocodes-public" + response.State.MaybePublicBucketWavAudioPath
			return audioDownloadPath

		default:
			log.Println("Failed to create audio file:", response.State.Status)
			return ""
		}
	}
}

func downloadAudioFile(path string) io.ReadCloser {
	resp, err := http.Get(path)
	if err != nil {
		log.Println("Error downloading audio file", err)
	}
	println("Content type is: ", resp.Header.Get("Content-Type"))
	return resp.Body
}
