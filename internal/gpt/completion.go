package gpt

import (
	"context"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"
)

func (g *GPTClient) CreateCompletion(ctx context.Context, prompt string) (string, error) {
	request := openai.CompletionRequest{
		Model:  openai.GPT3Ada,
		Prompt: prompt,
	}
	resp, err := g.api.CreateCompletion(ctx, request)
	if err != nil {
		return "", errors.Wrap(err, "creating completion")
	}
	return resp.Choices[0].Text, nil
}
