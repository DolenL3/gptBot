package gpt

import (
	"context"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"
)

// needs database to store users conversations
func (g *GPTClient) TalkToChatGPT(ctx context.Context, history []openai.ChatCompletionMessage, userID int64) (string, error) {

	request := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: history,
	}
	responce, err := g.api.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", errors.Wrap(err, "creating chat completion")
	}
	return responce.Choices[0].Message.Content, nil
}
