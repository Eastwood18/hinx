package hlog

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"hinx/utils"
	"os"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func (l *LogrusLogger) InfoF(format string, v ...interface{}) {

	l.logger.Infof(format, v...)
}

func (l *LogrusLogger) ErrorF(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

func (l *LogrusLogger) DebugF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	l.logger.Infof(format, v...)
}

func (l *LogrusLogger) PanicF(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusLogger) InfoFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusLogger) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusLogger) DebugFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusLogger) PanicFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func NewLogrusLogger() *LogrusLogger {
	out := utils.New("logs/app.log")
	out.SetCons(true)
	out.SetMaxSize(utils.SizeMiB / 100)
	//out.SetMaxAge(5)

	return &LogrusLogger{&logrus.Logger{
		Out:          out,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}}
}
