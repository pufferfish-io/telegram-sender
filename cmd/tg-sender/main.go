package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	cfg "tg-sender/internal/config"
	cfgModel "tg-sender/internal/config/model"
	"tg-sender/internal/messaging"
	"tg-sender/internal/processor"
	"tg-sender/internal/service"
)

func main() {
	log.Println("🚀 Starting tg-sender...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaConf, err := cfg.LoadSection[cfgModel.KafkaConfig]()
	if err != nil {
		log.Fatalf("❌ Failed to load Kafka config: %v", err)
	}

	tgConf, err := cfg.LoadSection[cfgModel.TelegramConfig]()
	if err != nil {
		log.Fatalf("❌ Failed to load Kafka config: %v", err)
	}

	if err != nil {
		log.Fatalf("❌ Failed to create S3 uploader: %v", err)
	}

	tgService := service.NewTelegramSenderService(tgConf.Token)

	tgMessagePreparer := processor.NewTgMessageSender(kafkaConf.TelegramMessageTopicName, tgService)
	handler := func(msg []byte) error {
		return tgMessagePreparer.Handle(ctx, msg)
	}

	messaging.Init(kafkaConf.BootstrapServersValue)

	consumer, err := messaging.NewConsumer(kafkaConf.BootstrapServersValue, kafkaConf.TelegramMessageGroupId, kafkaConf.TelegramMessageTopicName, handler)
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("❌ Consumer error: %v", err)
	}
}
