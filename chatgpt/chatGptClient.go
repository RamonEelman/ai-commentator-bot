package chatgpt

import (
	"context"
	"github.com/ayush6624/go-chatgpt"
	"log"
	"os"
)

var apiKey = os.Getenv("OPENAPI_KEY")

const systemMessage = "Respond as if you were Donald Trump after he became an presidential candidate. Incorporate trump-like sentences like build the wall and other trump-like sentences. \nThe score is the most important metric, don't mention specific amounts.Don't abbreviate, this output will be given to a text to speech model. Employ fake compliments, sarcasm, and comparisons to Donald Trump achievements.  Try to be nice to someone and mean to someone. Only use 22 words or less."
const userMessage = "Commentate on the following match result as donald trump only pick out 1 or 2 people to comment on: "

func StartChatGPTTask(inputChannel <-chan string) <-chan string {
	client, err := chatgpt.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	outputChannel := make(chan string)
	go func() {
		for {
			ctx := context.Background()
			matchSummary := <-inputChannel
			send, err := client.Send(ctx, &chatgpt.ChatCompletionRequest{
				Model: chatgpt.GPT4,
				Messages: []chatgpt.ChatMessage{
					{
						Role:    chatgpt.ChatGPTModelRoleSystem,
						Content: systemMessage,
					}, {
						Role:    chatgpt.ChatGPTModelRoleUser,
						Content: userMessage + matchSummary},
				},
				Temperature: 1.2,
			})
			if err != nil {
				log.Println("Could not Prompt chatgpt", err)
				continue
			}
			outputChannel <- send.Choices[0].Message.Content
		}
	}()
	return outputChannel
}
