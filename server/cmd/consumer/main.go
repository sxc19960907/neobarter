package main

import (
	"fmt"
	"log"

	"github.com/neobarter/server/internal/config"
	"github.com/neobarter/server/internal/pkg/es"
	"github.com/neobarter/server/internal/pkg/mq"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化 ES
	esClient, err := es.NewClient(cfg.Elasticsearch.Addresses)
	if err != nil {
		log.Fatalf("Failed to connect ES: %v", err)
	}

	if err := esClient.EnsureIndex(); err != nil {
		log.Fatalf("Failed to ensure index: %v", err)
	}

	// 初始化 MQ Consumer
	consumer, err := mq.NewConsumer(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to connect RabbitMQ: %v", err)
	}
	defer consumer.Close()

	log.Println("Consumer started, waiting for messages...")

	err = consumer.Consume(func(event mq.ItemEvent) {
		switch event.Type {
		case mq.EventCreate, mq.EventUpdate:
			if event.Data != nil {
				if err := esClient.IndexItem(event.ItemID, event.Data); err != nil {
					log.Printf("Failed to index item %d: %v", event.ItemID, err)
				} else {
					log.Printf("Indexed item %d", event.ItemID)
				}
			}
		case mq.EventDelete:
			if err := esClient.DeleteItem(event.ItemID); err != nil {
				log.Printf("Failed to delete item %d from index: %v", event.ItemID, err)
			} else {
				log.Printf("Deleted item %d from index", event.ItemID)
			}
		default:
			fmt.Printf("Unknown event type: %s\n", event.Type)
		}
	})

	if err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
