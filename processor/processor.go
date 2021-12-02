package processor

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

var (
	Errors       = errorx.NewNamespace("processor")
	CommonErrors = Errors.NewType("common")
)

type Config interface {
	MQZ2MAddress() string
	MQOpenHABAddress() string
	MQZ2MTopic() string
	MQOpenHABTopic() string
}

type Processor struct {
	cfg       Config
	log       *logrus.Entry
	z2mc, ohc mqtt.Client
}

func New(cfg Config, log *logrus.Logger) (*Processor, error) {
	p := &Processor{
		cfg: cfg,
		log: log.WithField("module", "processor"),
	}
	return p, nil
}

type Message struct {
	Topic   string
	Payload []byte
}

func (p *Processor) Process(ctx context.Context) error {
	//mqtt.DEBUG = p.log.WithField("mqtt_level", "debug")
	mqtt.WARN = p.log.WithField("mqtt_level", "warn")
	mqtt.CRITICAL = p.log.WithField("mqtt_level", "critical")
	mqtt.ERROR = p.log.WithField("mqtt_level", "error")
	z2mOpts := mqtt.NewClientOptions().
		AddBroker(p.cfg.MQZ2MAddress()).
		SetConnectRetryInterval(1 * time.Second)
	ohOpts := mqtt.NewClientOptions().
		AddBroker(p.cfg.MQOpenHABAddress()).
		SetConnectRetryInterval(1 * time.Second)
	t := time.NewTicker(1 * time.Second)
LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			if p.z2mc == nil || !p.z2mc.IsConnected() {
				p.log.Debugf("Connectiong to %s...", p.cfg.MQZ2MAddress())
				p.z2mc = mqtt.NewClient(z2mOpts)
				if token := p.z2mc.Connect(); token.Wait() && token.Error() != nil {
					p.log.Errorf("z2m connect: %+v", token.Error())
					continue
				}
			}
			if p.ohc == nil || !p.ohc.IsConnected() {
				p.log.Debugf("Connectiong to %s...", p.cfg.MQOpenHABAddress())
				p.ohc = mqtt.NewClient(ohOpts)
				if token := p.ohc.Connect(); token.Wait() && token.Error() != nil {
					p.log.Errorf("openHAB connect: %+v", token.Error())
					continue
				}
			}
			token := p.z2mc.Subscribe(p.cfg.MQZ2MTopic()+"/#", 0, func(client mqtt.Client, message mqtt.Message) {
				log := p.log.WithField("topic", message.Topic())
				var data DeviceData
				if err := json.Unmarshal(message.Payload(), &data); err != nil {
					log.Errorf("unmarshal: %+v", CommonErrors.WrapWithNoMessage(err))
					return
				}
				p.log.Debugf("Payload: %s", message.Payload())
				p.log.Debugf("Data: %+v", data)
				root := strings.TrimSuffix(p.cfg.MQOpenHABTopic()+strings.TrimPrefix(message.Topic(), p.cfg.MQZ2MTopic()), "/")
				msgs := data.Messages(root)
				go func() {
					for _, msg := range msgs {
						token := p.ohc.Publish(msg.Topic, 0, false, msg.Payload)
						if token.Wait() && token.Error() != nil {
							log.Errorf("publish error: %+v", token.Error())
						}
					}
				}()
			})
			if token.Wait() && token.Error() != nil {
				p.log.Errorf("subscribe error: %+v", token.Error())
				continue
			}
			break LOOP
		}
	}
	t.Stop()
	<-ctx.Done()
	return nil
}
