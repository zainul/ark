package xlog

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/zainul/zapbit"
	"go.uber.org/zap"
)

var (
	once        sync.Once
	wr          *zapbit.Writer
	err         error
	logger      *zap.Logger
	hostName    string
	serviceName string
)

// NewXLog is init streamer log
//func NewXLog(server, password, user, queueName, service string, port int) (*zap.Logger, *zapbit.Writer, error) {
func NewXLog(service string) (*zap.Logger, error) {
	once.Do(func() {
		logger, _ = zap.NewProduction(zap.WrapCore((&apmzap.Core{}).WrapCore))
		hostName, _ = os.Hostname()
		serviceName = service
	})

	return logger, nil

	// should trigger this in usage of xlog
	// Sync calls the underlying Core's Sync method, flushing any buffered log entries. Applications should take care to call Sync before exiting.
	// defer logger.Sync()
	// defer writer.Close()
}

// Info write to log info
func Info(ctx context.Context, msg string, value ...interface{}) {
	var (
		rqID string
		ok   bool
	)

	if rqID, ok = ctx.Value("request_id").(string); !ok {
		rqID = "NoRequestID"
	}

	var str string
	for _, val := range value {
		if val != nil {
			str = str + " " + val.(string)
		}
	}

	traceContextFields := apmzap.TraceContext(ctx)
	logger.With(traceContextFields...).Info(
		msg,
		zap.String("time", time.Now().Format(time.RFC3339)),
		zap.String("value", str),
		zap.String("request_id", rqID),
		zap.String("hostname", hostName),
		zap.String("service", serviceName),
		zap.ByteString("request_body", nil),
		zap.ByteString("response_body", nil),
	)
}

// Warning ...
func Warning(ctx context.Context, msg string, value ...interface{}) {

	var (
		rqID string
		ok   bool
	)

	if rqID, ok = ctx.Value("request_id").(string); !ok {
		rqID = "NoRequestID"
	}

	var str string
	for _, val := range value {
		if val != nil {
			str = str + " " + val.(string)
		}
	}

	traceContextFields := apmzap.TraceContext(ctx)
	logger.With(traceContextFields...).Warn(
		msg,
		zap.String("time", time.Now().Format(time.RFC3339)),
		zap.String("value", str),
		zap.String("request_id", rqID),
		zap.String("hostname", hostName),
		zap.String("service", serviceName),
		zap.ByteString("request_body", nil),
		zap.ByteString("response_body", nil),
	)
}

// Debug ...
func Debug(ctx context.Context, msg string, value ...interface{}) {
	traceContextFields := apmzap.TraceContext(ctx)
	bt, _ := json.Marshal(value)
	logger.With(traceContextFields...).Debug(
		msg,
		zap.String("time", time.Now().Format(time.RFC3339)),
		zap.ByteString("value", bt),
	)
}

// Error ...
func Error(ctx context.Context, msg string, err error) {

	var (
		rqID string
		ok   bool
	)

	if rqID, ok = ctx.Value("request_id").(string); !ok {
		rqID = "NoRequestID"
	}

	if err != nil {
		rqID = ctx.Value("request_id").(string)
		traceContextFields := apmzap.TraceContext(ctx)
		logger.With(traceContextFields...).Warn(
			msg+" :: "+err.Error(),
			zap.String("time", time.Now().Format(time.RFC3339)),
			zap.String("request_id", rqID),
			zap.String("hostname", hostName),
			zap.String("service", serviceName),
			zap.ByteString("request_body", nil),
			zap.ByteString("response_body", nil),
		)
	}

}

// Response ...
func Response(ctx context.Context, msg string, contractReq []byte, contractBody []byte) {

	rqID := ctx.Value("request_id").(string)

	traceContextFields := apmzap.TraceContext(ctx)
	logger.With(traceContextFields...).Warn(
		msg,
		zap.String("time", time.Now().Format(time.RFC3339)),
		zap.String("request_id", rqID),
		zap.String("hostname", hostName),
		zap.String("service", serviceName),
		zap.ByteString("request_body", contractReq),
		zap.ByteString("response_body", contractBody),
	)
}
