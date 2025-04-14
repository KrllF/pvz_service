package consumer

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/KrllF/pvz_service/internal/config"
)

type consumer struct {
	Consumer sarama.Consumer
	Conf     config.Config
}

// NewConsumer новый consumer
func NewConsumer(conf config.Config) (*consumer, error) {
	cons, err := newCons(strings.Split(conf.KafkaBrokers, ","))
	if err != nil {
		return nil, fmt.Errorf("newConsumer: %w", err)
	}

	return &consumer{Consumer: cons, Conf: conf}, nil
}

func newCons(brokerList []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	producer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (c *consumer) Close() {
	if err := c.Consumer.Close(); err != nil {
		log.Printf("c.Consumer.Close: %v", err)
	}
}
