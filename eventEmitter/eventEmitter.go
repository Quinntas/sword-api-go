package eventEmitter

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"sync"
	"time"
)

// ChannelManager handles multiple channels and their queues
type ChannelManager struct {
	conn     *amqp.Connection
	channels map[string]*ChannelInfo
	mu       sync.RWMutex
}

// ChannelInfo stores channel and its queues information
type ChannelInfo struct {
	channel *amqp.Channel
	queues  map[string]*QueueInfo
}

// QueueInfo stores queue information and its consumer (if any)
type QueueInfo struct {
	name     string
	messages <-chan amqp.Delivery
	consumer string
}

// NewChannelManager creates a new channel manager
func NewChannelManager(conn *amqp.Connection) *ChannelManager {
	return &ChannelManager{
		conn:     conn,
		channels: make(map[string]*ChannelInfo),
	}
}

// CreateChannel creates a new channel with the given name
func (cm *ChannelManager) CreateChannel(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.channels[name]; exists {
		return fmt.Errorf("channel %s already exists", name)
	}

	ch, err := cm.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	cm.channels[name] = &ChannelInfo{
		channel: ch,
		queues:  make(map[string]*QueueInfo),
	}

	return nil
}

// CreateQueue creates a new queue in the specified channel
func (cm *ChannelManager) CreateQueue(channelName, queueName string, durable, autoDelete, exclusive bool) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channelInfo, exists := cm.channels[channelName]
	if !exists {
		return fmt.Errorf("channel %s does not exist", channelName)
	}

	q, err := channelInfo.channel.QueueDeclare(
		queueName,  // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	channelInfo.queues[queueName] = &QueueInfo{
		name: q.Name,
	}

	return nil
}

// PublishMessage publishes a message to a specific queue
func (cm *ChannelManager) PublishMessage(channelName, queueName string, body []byte) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	channelInfo, exists := cm.channels[channelName]
	if !exists {
		return fmt.Errorf("channel %s does not exist", channelName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := channelInfo.channel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// ConsumeQueue starts consuming messages from a queue
func (cm *ChannelManager) ConsumeQueue(channelName, queueName string, handler func(amqp.Delivery)) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channelInfo, exists := cm.channels[channelName]
	if !exists {
		return fmt.Errorf("channel %s does not exist", channelName)
	}

	queueInfo, exists := channelInfo.queues[queueName]
	if !exists {
		return fmt.Errorf("queue %s does not exist in channel %s", queueName, channelName)
	}

	messages, err := channelInfo.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	queueInfo.messages = messages

	// Start consuming messages in a goroutine
	go func() {
		for msg := range messages {
			handler(msg)
		}
	}()

	return nil
}

func (cm *ChannelManager) CreateQueueWithConsumer(channelName, queueName string, consumer func(msg amqp.Delivery)) error {
	err := cm.CreateQueue(channelName, queueName, true, false, false)
	if err != nil {
		return err
	}
	err = cm.ConsumeQueue(channelName, queueName, consumer)
	if err != nil {
		return err
	}
	return nil
}

// CloseChannel closes a specific channel and all its queues
func (cm *ChannelManager) CloseChannel(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channelInfo, exists := cm.channels[name]
	if !exists {
		return fmt.Errorf("channel %s does not exist", name)
	}

	if err := channelInfo.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	delete(cm.channels, name)
	return nil
}

var Manager *ChannelManager

const (
	DefaultChannelName = "events"
)

func SetupEventEmitter() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}

	Manager = NewChannelManager(conn)
	if err := Manager.CreateChannel(DefaultChannelName); err != nil {
		log.Fatal(err)
	}
}
