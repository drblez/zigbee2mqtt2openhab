package logger

import (
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

var (
	Errors       = errorx.NewNamespace("logger")
	CommonErrors = Errors.NewType("common")
)

type Config interface {
	LoggerLevel() string
}

func New(cfg Config) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(cfg.LoggerLevel())
	if err != nil {
		return nil, CommonErrors.WrapWithNoMessage(err)
	}
	log := logrus.New()
	log.SetLevel(level)
	return log, nil
}
