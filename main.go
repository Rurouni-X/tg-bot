package main

import (
	"flag"
	"log"
	"tg-bot/clients/telegram"
	"tg-bot/consumer/event-consumer"
	eventTg "tg-bot/events/telegram"
	"tg-bot/storage/files"
)

const (
	tgHostBot = "api.telegram.org"
	storagePath = "storage"
	batchSize = 100
)

func main() {
	tgClient := telegram.NewClient(tgHostBot, mustToken())
	eventsProcessor := eventTg.NewProcessor(tgClient, files.NewFileStorage(storagePath))
	log.Print("service started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stoped", err)
	}
}

func mustToken() string {

	token := flag.String("token-bot", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
