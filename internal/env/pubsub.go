package env

type PubSubConfig struct{}

func (p *PubSubConfig) Host() string {
	return *pubsubHost
}

func (p *PubSubConfig) Project() string {
	return *pubsubProject
}
