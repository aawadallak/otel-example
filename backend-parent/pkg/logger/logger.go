package logger

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type (
	Logger interface {
		Infof(args ...any)
		Errof(args ...any)
	}

	Option func(l *logger)

	logger struct {
		logger *logrus.Logger
		tracer trace.Tracer
	}
)

func WithTrace(tracer trace.Tracer) Option {
	return func(l *logger) {
		l.tracer = tracer
	}
}

func (l *logger) Infof(args ...any) {
	l.logger.Printf("%+v", args...)
}

func (l *logger) Errof(args ...any) {
	l.logger.Printf("%+v", args...)
}

func New(opts ...Option) Logger {
	l := logger{logger: logrus.New()}

	for _, opt := range opts {
		opt(&l)
	}

	return &l
}
