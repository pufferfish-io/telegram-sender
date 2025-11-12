# telegram-sender

## Что делает

1. Подписывается на Kafka-топик `TOPIC_NAME_TG_REQUEST_MESSAGE`, входящий поток подготовленных `SendMessageRequest` (см. `internal/contract`).
2. Каждый сырое сообщение десериализуется в структуру с `chat_id`, `text` и опциональными параметрами, затем преобразуется в HTTP-запрос к Telegram API.
3. Отправляет POST на `https://api.telegram.org/bot<TOKEN>/sendMessage` и логирует исходные данные и код ответа.
4. При ошибках клиента или Telegram API возвращает ненулевой код, чтобы `runConsumerSupervisor` перезапустил consumer.

## Запуск

1. Экспортируйте нужные переменные (`set -a && source .env && set +a`) или задавайте вручную.
2. Соберите и запустите локально:
   ```bash
   go run ./cmd/tg-sender
   ```
3. Или соберите Docker-образ и пробросьте конфигурацию:
   ```bash
   docker build -t telegram-sender .
   docker run --rm -e ... telegram-sender
   ```

## Переменные окружения

Все переменные обязательны, кроме SASL-полей, если Kafka не требует аутентификации.

- `KAFKA_BOOTSTRAP_SERVERS_VALUE` — список брокеров (`host:port[,host:port]`).
- `KAFKA_TOPIC_NAME_TG_REQUEST_MESSAGE` — топик, откуда читаются `SendMessageRequest` для Telegram.
- `KAFKA_GROUP_ID_TELEGRAM_SENDER` — consumer group для масштабирования обработки ответов.
- `KAFKA_CLIENT_ID_TELEGRAM_SENDER` — общий client id (продюсер+консьюмер), отображается в метриках Sarama.
- `KAFKA_SASL_USERNAME` и `KAFKA_SASL_PASSWORD` — используйте для SASL/SCRAM, можно оставить пустыми при открытом кластере.
- `TELEGRAM_TOKEN` — токен бота в формате `1234:abcd...`, используется для формирования URL `https://api.telegram.org/bot<TOKEN>`.

## Примечания

- Сам HTTP-запрос формируется в `internal/processor/tg-message-sender.go`, `SendMessageRequest` маршалится через `encoding/json`.
- При ответе Telegram со статусом ≥300 или когда тело не парсится, ошибка возвращается, и `runConsumerSupervisor` запускает consumer снова (логика в `cmd/tg-sender/main.go`).
- Сообщения отправляются синхронно, поэтому скорость ограничена задержками Telegram API; обрабатываются по одному.
