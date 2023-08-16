package gpt

import openai "github.com/sashabaranov/go-openai"

type GPTClient struct {
	api   openai.Client
	token string
}

func New(token string) *GPTClient {
	return &GPTClient{api: *openai.NewClient(token), token: token}
}
