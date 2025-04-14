package producer

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

type producer struct {
	SyncProducer sarama.SyncProducer
}

// NewProducer конструктор producer
func NewProducer(brokerAddress string) (*producer, error) {
	prod, err := newSyncProducer(strings.Split(brokerAddress, ","))
	if err != nil {
		return nil, fmt.Errorf("newSyncProducer: %w", err)
	}

	return &producer{SyncProducer: prod}, nil
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (p *producer) Close() {
	if err := p.SyncProducer.Close(); err != nil {
		log.Println("p.SyncProducer.Close()")
	}
}
