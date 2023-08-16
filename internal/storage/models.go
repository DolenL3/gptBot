package storage

import "github.com/sashabaranov/go-openai"

type ChatCompletionMessage struct {
	Role         string               `bson:"role"`
	Content      string               `bson:"content"`
	Name         string               `bson:"name,omitempty"`
	FunctionCall *openai.FunctionCall `bson:"function_call,omitempty"`
}

type History struct {
	Chat []ChatCompletionMessage `bson:"chat,omitempty"`
}

type User struct {
	TgID int64                   `bson:"tgID"`
	Chat []ChatCompletionMessage `bson:"chat,omitempty"`
}
