package producer

import (
	"fmt"

	"github.com/IBM/sarama"
)

func (p *producer) SendMessage(message []byte, topicname string) error {
	msg := &sarama.ProducerMessage{
		Topic: topicname,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := p.SyncProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("p.SyncProducer.SendMessage: %w", err)
	}

	return nil
}
