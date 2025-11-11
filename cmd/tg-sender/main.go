package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cfg "tg-sender/internal/config"
	"tg-sender/internal/logger"
	"tg-sender/internal/messaging"
	"tg-sender/internal/processor"
)

func main() {
	lg, cleanup := logger.NewZapLogger()
	defer cleanup()
	lg.Info("🚀 Starting tg-sender...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lg.Info("🚀 Starting tg-sender…")

	conf, err := cfg.Load()
	if err != nil {
		lg.Error("❌ Failed to load config: %v", err)
		os.Exit(1)
	}

	sender := processor.NewTgMessageSender(processor.Option{
		Token:      conf.Telegram.Token,
		ApiBase:    "https://api.telegram.org/bot" + conf.Telegram.Token,
		HttpClient: http.DefaultClient,
		Logger:     lg,
	})

	consumer, err := messaging.NewKafkaConsumer(messaging.ConsumerOption{
		Logger:       lg,
		Broker:       conf.Kafka.BootstrapServersValue,
		GroupID:      conf.Kafka.ResponseMessageGroupID, // используем имеющийся group id
		Topics:       []string{conf.Kafka.TelegramMessageTopicName},
		Handler:      sender,
		SaslUsername: conf.Kafka.SaslUsername,
		SaslPassword: conf.Kafka.SaslPassword,
		ClientID:     conf.Kafka.ClientID,
		Context:      ctx,
	})
	if err != nil {
		lg.Error("❌ Failed to create consumer: %v", err)
		os.Exit(1)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	if err := consumer.Start(ctx); err != nil {
		lg.Error("❌ Consumer error: %v", err)
		os.Exit(1)
	}
}
