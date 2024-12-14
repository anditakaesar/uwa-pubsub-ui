package xpubsub

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cloud.google.com/go/pubsub"
)

type TopicItem struct {
	ID   string
	Name string
}

type SubscriptionItem struct {
	ID      string
	Name    string
	TopicID string
}

type PubSubClient struct {
	client *pubsub.Client
	config ConfigInterface
}

func NewPubSubClient(config ConfigInterface) *PubSubClient {
	return &PubSubClient{
		config: config,
	}
}

type ConfigInterface interface {
	Host() string
	Project() string
}

type Interface interface {
	GetTopics(ctx context.Context) ([]TopicItem, error)
	GetSubscriptionsByTopic(ctx context.Context, topic string) ([]SubscriptionItem, error)
	CreateTopic(ctx context.Context, id string) error
	GetSubscriptions(ctx context.Context) ([]SubscriptionItem, error)
	CreateSubscription(ctx context.Context, id string, topicID string) error
	DeleteTopic(ctx context.Context, id string) error
	DeleteSubscription(ctx context.Context, id string) error
	CreateTopicAndSubscriptions(ctx context.Context, topicSpecs []TopicSpec) error
	// PublishMessage(ctx context.Context, topicID string, data []byte) error
	// SubscribeMessage(ctx context.Context, subscriptionID string) (<-chan *pubsub.Message, error)
	// AcknowledgeMessage(ctx context.Context, msg *pubsub.Message) error
}

func (c *PubSubClient) getClient(ctx context.Context) (func() error, error) {
	client, err := pubsub.NewClient(ctx, c.config.Project())
	if err != nil {
		return nil, err
	}
	c.client = client
	return client.Close, nil
}

func (c *PubSubClient) GetTopics(ctx context.Context) ([]TopicItem, error) {
	close, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer close()

	topicsIter := c.client.Topics(ctx)
	var topics []TopicItem

	for {
		topic, err := topicsIter.Next()
		if err != nil {
			break
		}

		item := TopicItem{
			ID:   topic.ID(),
			Name: topic.String(),
		}

		topics = append(topics, item)
	}

	return topics, nil
}

func (c *PubSubClient) GetSubscriptionsByTopic(ctx context.Context, topic string) ([]SubscriptionItem, error) {
	close, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer close()

	subsIter := c.client.Subscriptions(ctx)
	var subs []SubscriptionItem
	for {
		sub, err := subsIter.Next()
		if err != nil {
			break
		}
		cfg, err := sub.Config(ctx)
		if err != nil {
			break
		}

		exist, err := cfg.Topic.Exists(ctx)
		if err != nil {
			break
		}

		if !exist {
			continue
		}

		if cfg.Topic.ID() == topic {
			item := SubscriptionItem{
				ID:      sub.ID(),
				Name:    sub.String(),
				TopicID: topic,
			}
			subs = append(subs, item)
		}
	}
	return subs, nil
}

func (c *PubSubClient) CreateTopic(ctx context.Context, topic string) error {
	close, err := c.getClient(ctx)
	if err != nil {
		return err
	}
	defer close()
	_, err = c.client.CreateTopic(ctx, topic)
	return err
}

func (c *PubSubClient) CreateSubscription(ctx context.Context, subscription string, topic string) error {
	close, err := c.getClient(ctx)
	if err != nil {
		return err
	}
	defer close()
	_, err = c.client.CreateSubscription(ctx, subscription, pubsub.SubscriptionConfig{
		Topic:       c.client.Topic(topic),
		AckDeadline: 20 * time.Second,
	})
	return err
}

func (c *PubSubClient) GetSubscriptions(ctx context.Context) ([]SubscriptionItem, error) {
	close, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer close()
	subsIter := c.client.Subscriptions(ctx)
	var subs []SubscriptionItem
	for {
		sub, err := subsIter.Next()
		if err != nil {
			break
		}

		cfg, err := sub.Config(ctx)
		if err != nil {
			return nil, err
		}

		item := SubscriptionItem{
			ID:   sub.ID(),
			Name: sub.String(),
		}

		exist, err := cfg.Topic.Exists(ctx)
		if err != nil {
			return nil, err
		}
		if exist {
			item.TopicID = cfg.Topic.String()
		} else {
			item.TopicID = "unknown (deleted)"
		}

		subs = append(subs, item)
	}
	return subs, nil
}

func (c PubSubClient) getSubscriptionsByTopicForDelete(ctx context.Context, topic string) ([]SubscriptionItem, error) {
	subsIter := c.client.Subscriptions(ctx)
	var subs []SubscriptionItem
	for {
		sub, err := subsIter.Next()
		if err != nil {
			break
		}
		cfg, err := sub.Config(ctx)
		if err != nil {
			break
		}

		exist, err := cfg.Topic.Exists(ctx)
		if err != nil {
			break
		}

		if !exist {
			continue
		}

		if cfg.Topic.ID() == topic {
			item := SubscriptionItem{
				ID:      sub.ID(),
				Name:    sub.String(),
				TopicID: topic,
			}
			subs = append(subs, item)
		}
	}
	return subs, nil
}

func (c *PubSubClient) DeleteTopic(ctx context.Context, topic string) error {
	close, err := c.getClient(ctx)
	if err != nil {
		return err
	}
	defer close()

	subs, err := c.getSubscriptionsByTopicForDelete(ctx, topic)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		go func(subName string) {
			deleteSubErr := c.DeleteSubscription(context.Background(), subName)
			if deleteSubErr != nil {
				slog.Error("failed to delete subscription", "error", deleteSubErr)
			}
		}(sub.ID)
	}

	return c.client.Topic(topic).Delete(ctx)
}

func (c *PubSubClient) DeleteSubscription(ctx context.Context, subsName string) error {
	close, err := c.getClient(ctx)
	if err != nil {
		return err
	}
	defer close()
	subs := c.client.Subscription(subsName)

	return subs.Delete(ctx)
}

func (c *PubSubClient) PublishMessage(ctx context.Context, topic string, msg string) error {
	close, err := c.getClient(ctx)
	if err != nil {
		return err
	}
	defer close()
	t := c.client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	_, err = result.Get(ctx)
	return err
}

func (c *PubSubClient) CreateTopicAndSubscriptions(ctx context.Context, topicSpecs []TopicSpec) error {
	close, err := c.getClient(ctx)
	if err != nil {
		return err
	}
	defer close()
	start := time.Now()
	defer func(startTime time.Time) {
		endTime := time.Since(startTime)
		slog.Info(fmt.Sprintf("==== Creating Topics Done: %s ====", endTime))
	}(start)
	slog.Info("==== Creating Topics ====")

	for _, topicSpec := range topicSpecs {
		slog.Info(fmt.Sprintf("creating topic: %s with %d subscribers", topicSpec.Name, len(topicSpec.Subscriptions)))
		topic := c.client.Topic(topicSpec.Name)
		exists, err := topic.Exists(ctx)
		if err != nil {
			return err
		}
		if !exists {
			slog.Info("topic not exist, creating...")
			topic, err = c.client.CreateTopic(ctx, topicSpec.Name)
			if err != nil {
				return err
			}
		}

		for _, subsSpec := range topicSpec.Subscriptions {
			slog.Info(fmt.Sprintf("creating subscription: %s", subsSpec))
			subs := c.client.Subscription(subsSpec)
			exists, err := subs.Exists(ctx)
			if err != nil {
				return err
			}
			if !exists {
				slog.Info("subscription not exist, creating...")
				_, err = c.client.CreateSubscription(ctx, subsSpec, pubsub.SubscriptionConfig{
					Topic:       topic,
					AckDeadline: 20 * time.Second,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
