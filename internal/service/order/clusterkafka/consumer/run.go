package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

func (c *consumer) Run(ctx context.Context) error {
	msgCnt := 0

	consumer, err := c.Consumer.ConsumePartition(c.Conf.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("worker.ConsumePartition: %w", err)
	}

	log.Println("Consumer запущен")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				msgCnt++
				fmt.Println(string(msg.Value))
			}
		}
	}()

	wg.Wait()

	return nil
}
