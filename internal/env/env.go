package env

import (
	"os"
	"sync"
)

var (
	pubsubHost    *string
	pubsubProject *string
	once          sync.Once
)

func InitEnv() {
	once.Do(func() {
		pubsubHost = new(string)
		*pubsubHost = os.Getenv("PUBSUB_EMULATOR_HOST")

		pubsubProject = new(string)
		*pubsubProject = os.Getenv("ProjectID")
	})
}

func PubSubInitFile() string {
	return os.Getenv("PubSubInitFile")
}

func InitPubSubFromFile() bool {
	return os.Getenv("InitPubSubFromFile") == "true"
}
