package broker

import (
	"encoding/json"
	"fmt"

	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/config"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/models"
	"github.com/streadway/amqp"
)

const (
	TOPIC string = "topic"
)

type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewBroker(cfg *config.ServiceConfig) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.RabbitMQConfig.URI)
	if err != nil {
		return nil, fmt.Errorf("dial rabbit: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // nolint:errcheck, gosec
		return nil, fmt.Errorf("open conn channel: %w", err)
	}

	if err := ch.ExchangeDeclare(cfg.RabbitMQConfig.Exchange, TOPIC,
		true, false, false, false, nil); err != nil {
		conn.Close() // nolint:errcheck, gosec
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	return &RabbitMQ{conn: conn, channel: ch, exchange: cfg.RabbitMQConfig.Exchange}, nil
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		r.conn.Close() // nolint:errcheck, gosec
		return fmt.Errorf("close rabbit channel: %w", err)
	}

	return r.conn.Close()
}

func (r *RabbitMQ) PublishEvent(order *models.Order) error {
	body, err := json.Marshal(Message{
		ID:            order.ID.String(),
		Status:        string(order.Status),
		TotalPrice:    order.TotalPrice,
		TotalQuantity: order.TotalQuantity,
	})
	if err != nil {
		return fmt.Errorf("marshal order to message: %w", err)
	}

	queue, err := r.channel.QueueDeclare("orders_queue", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	if err := r.channel.QueueBind(queue.Name, "orders.created", r.exchange, false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	// some shi must be from config.yml but now i want to sleep a lot
	if err := r.channel.Publish(r.exchange, "orders.created", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}

type Message struct {
	ID            string  `json:"id"`
	Status        string  `json:"status"`
	TotalPrice    float64 `json:"total_price"`
	TotalQuantity int     `json:"total_quantity"`
}
