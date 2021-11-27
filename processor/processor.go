package processor

import (
	"context"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type Processor struct {
}

func New(log *logrus.Logger) (*Processor, error) {
	return &Processor{}, nil
}

func (p *Processor) Process(ctx context.Context) error {
	mqtt.ConnectMQTT()
}
