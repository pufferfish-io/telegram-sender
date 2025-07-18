package model

type KafkaConfig struct {
	BootstrapServersValue    string `yaml:"bootstrap_servers_value"`
	TelegramMessageTopicName string `yaml:"telegram_message_topic_name"`
	TelegramMessageGroupId   string `yaml:"telegram_message_group_id"`
}

func (KafkaConfig) SectionName() string {
	return "kafka"
}
