package service

import (
	"context"

	"zigbee2mqtt2openhab/config"
	"zigbee2mqtt2openhab/logger"
	"zigbee2mqtt2openhab/processor"

	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func (svc *Service) makeContainer() (*dig.Container, error) {
	c := dig.New()
	var err error
	add := func(constructor interface{}) {
		if err != nil {
			return
		}
		err = c.Provide(constructor)
	}
	add(func() context.Context {
		return svc.ctx
	})
	add(func(cfg *config.Config) (*logrus.Logger, error) {
		return logger.New(cfg)
	})
	add(func() (*config.Config, error) {
		return config.New()
	})
	add(func(cfg *config.Config, log *logrus.Logger) (*processor.Processor, error) {
		return processor.New(log)
	})
	if err != nil {
		return nil, CommonErrors.WrapWithNoMessage(err)
	}
	return c, nil
}
