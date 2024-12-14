package xpubsub

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TopicSpec struct {
	Name          string
	Subscriptions []string
}

type TopicWrapper struct {
	Topics map[string][]string `yaml:"topics"`
}

func LoadTopic(file string) ([]TopicSpec, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var topics TopicWrapper
	err = yaml.Unmarshal(data, &topics)
	if err != nil {
		return nil, err
	}

	var topicSpecs []TopicSpec
	for name, subscribers := range topics.Topics {
		topicSpecs = append(topicSpecs, TopicSpec{
			Name:          name,
			Subscriptions: subscribers,
		})
	}
	return topicSpecs, nil
}
