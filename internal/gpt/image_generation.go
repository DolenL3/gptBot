package gpt

import (
	"context"
	"encoding/base64"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"
)

func (g *GPTClient) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	request := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
	b64JSONImg, err := g.api.CreateImage(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "creating image")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(b64JSONImg.Data[0].B64JSON)
	if err != nil {
		return nil, errors.Wrap(err, "decoding into image")
	}

	return imgBytes, nil
}
