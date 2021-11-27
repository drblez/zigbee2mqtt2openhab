package service

import (
	"context"
	"flag"
	"strings"

	"zigbee2mqtt2openhab/processor"

	"github.com/joomcode/errorx"
	"github.com/kardianos/service"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
)

var (
	Errors       = errorx.NewNamespace("service")
	CommonErrors = Errors.NewType("common")
)

var (
	command = flag.String("cmd", "", "Service command: "+
		strings.Join(service.ControlAction[:], ","))
)

type Service struct {
	svc    service.Service
	logger service.Logger

	ctx    context.Context
	cancel context.CancelFunc

	g *errgroup.Group
}

func New() (*Service, error) {
	s := &Service{}
	cfg := &service.Config{
		Name:        "zigbee2mqtt2openhab",
		DisplayName: "ZigBee => MQTT => OpenHAB",
		Description: "ZigBee to MQTT to OpenHAB",
	}
	var err error
	s.svc, err = service.New(s, cfg)
	if err != nil {
		return nil, CommonErrors.WrapWithNoMessage(err)
	}
	s.logger, err = s.svc.Logger(nil)
	if err != nil {
		return nil, CommonErrors.WrapWithNoMessage(err)
	}
	return s, nil
}

func (svc *Service) Run() error {
	switch *command {
	case "":
	case service.ControlAction[0]:
	case service.ControlAction[1]:
	case service.ControlAction[2]:
	case service.ControlAction[3]:
	case service.ControlAction[4]:
	default:
		return CommonErrors.New("unknown command %q", *command)
	}
	if *command != "" {
		if err := service.Control(svc.svc, *command); err != nil {
			return CommonErrors.WrapWithNoMessage(err)
		}
		return nil
	}
	if err := svc.svc.Run(); err != nil {
		return CommonErrors.WrapWithNoMessage(err)
	}
	return nil
}

func (svc *Service) run(c *dig.Container) error {
	err := c.Invoke(func(ctx context.Context, p *processor.Processor) error {
		return p.Process(ctx)
	})
	if err != nil {
		return CommonErrors.WrapWithNoMessage(err)
	}
	return nil
}

func (svc *Service) Start(_ service.Service) error {
	_ = svc.logger.Info("Starting service...")
	svc.ctx, svc.cancel = context.WithCancel(context.Background())
	svc.g = &errgroup.Group{}
	c, err := svc.makeContainer()
	if err != nil {
		_ = svc.logger.Errorf("Starting service: %+v", err)
		return err
	}
	svc.g.Go(func() error {
		err := svc.run(c)
		if err != nil {
			panic(err)
		}
		return nil
	})
	return nil
}

func (svc *Service) Stop(_ service.Service) error {
	_ = svc.logger.Info("Stopping service...")
	svc.cancel()
	if err := svc.g.Wait(); err != nil {
		return CommonErrors.WrapWithNoMessage(err)
	}
	return nil
}
