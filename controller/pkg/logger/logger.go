package logger

import (
	formatter "github.com/fabienm/go-logrus-formatters"
	"github.com/sirupsen/logrus"
)

type Log struct {
	*logrus.Logger
}

func NewLogger(logLvl logrus.Level) *Log {
	l := &Log{logrus.New()}
	l.SetLevel(logLvl)
	gelfFormat := formatter.NewGelf("Controller-Service")
	l.SetFormatter(gelfFormat)
	return l
}
