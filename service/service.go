package service

import (
	"context"
	"flag"
	"github.com/joomcode/errorx"
	"github.com/kardianos/service"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
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

func (svc *Service) run() error {
	ctx := svc.ctx
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		time.Sleep(1 * time.Second)
	}
}

func (svc *Service) Start(_ service.Service) error {
	_ = svc.logger.Info("Starting service...")
	svc.ctx, svc.cancel = context.WithCancel(context.Background())
	svc.g = &errgroup.Group{}
	svc.g.Go(svc.run)
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
