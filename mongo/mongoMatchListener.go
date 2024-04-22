package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var MONGO_URI = os.Getenv("MONGO_URI")

type MatchSummaryEntity struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Summary string             `bson:"summary,omitempty"`
}

func StartChangeListener() <-chan string {
	if MONGO_URI == "" {
		log.Fatal("MONGO_URI was not set")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		panic(err)
	}
	pipeline := mongo.Pipeline{}
	opts := options.ChangeStream().SetFullDocument(options.Required)

	changestream, err := client.Database("ivanbot").Collection("match-summary").Watch(context.TODO(), pipeline, opts)
	if err != nil {
		log.Println("Error while creating change stream: ", err)
	}
	channel := make(chan string)
	go watchChangeStream(context.TODO(), changestream, channel)
	return channel
}
func watchChangeStream(ctx context.Context, changeStream *mongo.ChangeStream, channel chan<- string) {
	defer changeStream.Close(ctx)
	log.Println("Watching change stream")

	for changeStream.Next(ctx) {
		var event struct {
			FullDocument MatchSummaryEntity `bson:"fullDocument"`
		}
		if err := changeStream.Decode(&event); err != nil {
			log.Println("Error while decoding change doc: ", err)
		}
		channel <- event.FullDocument.Summary
		log.Println("found new match summary")
		fmt.Printf("Received MatchSummaryEntity: %+v\n", event.FullDocument)
	}
}
