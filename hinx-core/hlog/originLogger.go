package hlog

import (
	"context"
	"log"
)

type originLogger struct{}

func (o *originLogger) PanicF(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

func (o *originLogger) PanicFX(ctx context.Context, format string, v ...interface{}) {
	log.Println(ctx)
	log.Panicf(format, v...)
}

func (o *originLogger) InfoF(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (o *originLogger) ErrorF(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (o *originLogger) DebugF(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (o *originLogger) InfoFX(ctx context.Context, format string, v ...interface{}) {
	log.Println(ctx)
	log.Printf(format, v...)
}

func (o *originLogger) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	log.Println(ctx)
	log.Printf(format, v...)
}

func (o *originLogger) DebugFX(ctx context.Context, format string, v ...interface{}) {
	log.Println(ctx)
	log.Printf(format, v...)
}
