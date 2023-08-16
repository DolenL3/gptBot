package bot

import (
	"context"
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"

	"gptbot/internal/gpt"
	"gptbot/internal/storage"
)

type Bot struct {
	api *tgbotapi.BotAPI
	gpt *gpt.GPTClient
	db  *storage.Storage
}

func New(api *tgbotapi.BotAPI, gpt *gpt.GPTClient, db *storage.Storage) *Bot {
	return &Bot{api: api, gpt: gpt, db: db}
}

func (b *Bot) Start(ctx context.Context) error {
	updates := b.api.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})

	wg := &sync.WaitGroup{}

	for {
		select {
		case update := <-updates:
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := b.handleUpdate(ctx, update)
				if err != nil {
					log.Printf("can't handle update: %v\n", err)
				}
			}()
		case <-ctx.Done():
			wg.Wait()
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	err := b.handleMessage(ctx, update.Message)
	if err != nil {
		_, errSend := b.api.Send(tgbotapi.NewMessage(update.Message.From.ID, _internalErrorReply))
		if errSend != nil {
			log.Printf("can't send internal error reply: %v", err)
		}
		return errors.Wrap(err, "handle message")
	}
	return nil
}

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) error {

	switch msg.Command() {
	case _chatGPT:
		prompt, _ := strings.CutPrefix(msg.Text, "/"+_chatGPT)

		b.api.Send(tgbotapi.NewMessage(msg.From.ID, "запрос передан на сохранение"))
		history, err := b.db.AddHistory(ctx, msg.From.ID, storage.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: prompt})
		if err != nil {
			return errors.Wrap(err, "handling message")
		}

		b.api.Send(tgbotapi.NewMessage(msg.From.ID, "запрос передан в gpt api"))
		responce, err := b.gpt.TalkToChatGPT(ctx, history, msg.From.ID)
		if err != nil {
			return errors.Wrap(err, "creating bot responce")
		}

		b.api.Send(tgbotapi.NewMessage(msg.From.ID, "ответ передан на сохранение"))
		_, err = b.db.AddHistory(ctx, msg.From.ID, storage.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: responce})
		if err != nil {
			b.api.Send(tgbotapi.NewMessage(msg.From.ID, "ошибка сохранения ответа"))
		}

		b.api.Send(tgbotapi.NewMessage(msg.From.ID, responce))

	case _chatCompletion:
		prompt, _ := strings.CutPrefix(msg.Text, "/"+_chatCompletion)

		b.api.Send(tgbotapi.NewMessage(msg.From.ID, "запрос передан в gpt api"))
		responce, err := b.gpt.CreateCompletion(ctx, prompt)
		if err != nil {
			return errors.Wrap(err, "creating chat completion")
		}

		b.api.Send(tgbotapi.NewMessage(msg.From.ID, responce))

	case _imageGeneration:
		prompt, _ := strings.CutPrefix(msg.Text, "/"+_imageGeneration)

		imgBytes, err := b.gpt.GenerateImage(ctx, prompt)
		if err != nil {
			return errors.Wrap(err, "creating image")
		}

		img := tgbotapi.NewPhoto(msg.From.ID, tgbotapi.FileBytes{Name: "generated image", Bytes: imgBytes})
		b.api.Send(img)
	}

	return nil
}
