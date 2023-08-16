package main

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gptbot/internal/bot"
	"gptbot/internal/gpt"
	"gptbot/internal/storage"
)

func main() {
	fmt.Println("started")
	err := run()
	if err != nil {
		fmt.Printf("run: %v\n", err)
		os.Exit(1)
	}

}

func run() error {
	err := godotenv.Load()
	if err != nil {
		return errors.Wrap(err, "loading .env")
	}

	tgToken := os.Getenv("TELEGRAM_BOT_KEY")
	if tgToken == "" {
		return errors.New("no telegram api key found in .env file")
	}
	gptToken := os.Getenv("GPT_API_KEY")
	if gptToken == "" {
		return errors.New("no gpt api key found in .env file")
	}

	api, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return errors.Wrap(err, "new bot api")
	}

	ctx := context.TODO()

	gpt := gpt.New(gptToken)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		return errors.Wrap(err, "connect db")
	}

	db := storage.New(mongoClient)

	bot := bot.New(api, gpt, db)

	err = bot.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "start bot")
	}

	return nil
}
