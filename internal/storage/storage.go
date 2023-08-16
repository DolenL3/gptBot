package storage

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage is storage implementation via mongodb.
type Storage struct {
	client *mongo.Client
}

// New return Storage implemented via mongodb.
func New(client *mongo.Client) *Storage {
	return &Storage{client: client}
}

// AddHistory adds users prompt to conversation history and returns all history.
func (s *Storage) AddHistory(ctx context.Context, tgID int64, prompt ChatCompletionMessage) ([]openai.ChatCompletionMessage, error) {
	users := s.client.Database("GPTData").Collection("users")

	filter := bson.D{primitive.E{Key: "tgID", Value: tgID}}
	update := bson.D{primitive.E{Key: "$push", Value: bson.D{
		primitive.E{Key: "chat", Value: prompt}}}}

	updateResult, err := users.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "updating database")
	}

	if updateResult.MatchedCount == 0 {
		log.Printf("user %d was not found, adding to db and retrying...", tgID)
		users.InsertOne(ctx, User{TgID: tgID})
		return s.AddHistory(ctx, tgID, prompt)
	}

	var history History
	findUserFilter := bson.D{{Key: "tgID", Value: tgID}}
	findHistoryOption := bson.D{{Key: "chat", Value: 1}, {Key: "_id", Value: 0}}
	opts := options.FindOne().SetProjection(findHistoryOption)

	a := users.FindOne(ctx, findUserFilter, opts)
	err = a.Decode(&history)
	if err != nil {
		return nil, errors.Wrap(err, "decoding find result")
	}

	// convert history to slice of openai ChatCompletionMessage
	var openaiHistory []openai.ChatCompletionMessage
	for _, message := range history.Chat {
		openaiHistory = append(openaiHistory, storageToOpenaiMessage(message))
	}

	return openaiHistory, nil
}

func storageToOpenaiMessage(msg ChatCompletionMessage) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:         msg.Role,
		Content:      msg.Content,
		Name:         msg.Name,
		FunctionCall: msg.FunctionCall,
	}
}
