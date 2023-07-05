package hlog

import (
	"context"
	"hinx/hinx-core/hiface"
)

var hLogInstance hiface.ILogger = new(hinxDefaultLog)

type hinxDefaultLog struct{}

func (h *hinxDefaultLog) PanicF(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (h *hinxDefaultLog) PanicFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (h *hinxDefaultLog) InfoF(format string, v ...interface{}) {
	StdHinxLog.Infof(format, v...)
}

func (h *hinxDefaultLog) ErrorF(format string, v ...interface{}) {
	StdHinxLog.Errorf(format, v...)
}

func (h *hinxDefaultLog) DebugF(format string, v ...interface{}) {
	StdHinxLog.Debugf(format, v...)
}

func (h *hinxDefaultLog) InfoFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (h *hinxDefaultLog) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (h *hinxDefaultLog) DebugFX(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func SetLogger(newLog hiface.ILogger) {
	hLogInstance = newLog
}

func Ins() hiface.ILogger {
	return hLogInstance
}
