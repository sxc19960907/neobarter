package mq

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "neobarter"
	ItemQueue    = "item_index_sync"
	ItemRouteKey = "item.sync"
)

// EventType 事件类型
type EventType string

const (
	EventCreate EventType = "create"
	EventUpdate EventType = "update"
	EventDelete EventType = "delete"
)

// ItemEvent 物品变更事件
type ItemEvent struct {
	Type   EventType              `json:"type"`
	ItemID int64                  `json:"item_id"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

// Publisher 消息发布者
type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明交换机
	err = ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Println("RabbitMQ publisher connected")
	return &Publisher{conn: conn, channel: ch}, nil
}

func (p *Publisher) PublishItemEvent(event ItemEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		ExchangeName,
		ItemRouteKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) Close() {
	p.channel.Close()
	p.conn.Close()
}

// Consumer 消息消费者
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明交换机和队列
	err = ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(ItemQueue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(ItemQueue, ItemRouteKey, ExchangeName, false, nil)
	if err != nil {
		return nil, err
	}

	log.Println("RabbitMQ consumer connected")
	return &Consumer{conn: conn, channel: ch}, nil
}

func (c *Consumer) Consume(handler func(ItemEvent)) error {
	msgs, err := c.channel.Consume(ItemQueue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var event ItemEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}
		handler(event)
	}
	return nil
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
