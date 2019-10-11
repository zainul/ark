package xlog

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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
func NewXLog(server, password, user, queueName, service string, port int) (*zap.Logger, *zapbit.Writer, error) {
	once.Do(func() {
		if strings.ToLower(os.Getenv("env")) == strings.ToLower("production") {
			wr, err = zapbit.NewWriter(zapbit.RabbitMQConfig{
				Address:  server,
				Password: password,
				User:     user,
				Port:     port,
			}, queueName)

			if err != nil {
				fmt.Println("error connect to rabbit mq streamer", err)
				return
			}

			log.Println("this message for log Production Only")
			logger = zap.New(wr.GetCore())

		} else {
			logger = zap.NewExample() // or NewProduction, or NewDevelopment
		}
		hostName, _ = os.Hostname()
		serviceName = service
	})

	return logger, wr, nil

	// should trigger this in usage of xlog
	// Sync calls the underlying Core's Sync method, flushing any buffered log entries. Applications should take care to call Sync before exiting.
	// defer logger.Sync()
	// defer writer.Close()
}

// Info write to log info
func Info(ctx context.Context, msg string, value ...interface{}) {

	rqID := ctx.Value("request_id").(string)

	var str string
	for _, val := range value {
		if val != nil {
			str = str + " " + val.(string)
		}
	}

	logger.Info(
		msg,
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

	rqID := ctx.Value("request_id").(string)

	var str string
	for _, val := range value {
		if val != nil {
			str = str + " " + val.(string)
		}
	}

	logger.Warn(
		msg,
		zap.String("value", str),
		zap.String("request_id", rqID),
		zap.String("hostname", hostName),
		zap.String("service", serviceName),
		zap.ByteString("request_body", nil),
		zap.ByteString("response_body", nil),
	)
}

// Debug ...
func Debug(msg string, value ...interface{}) {
	logger.Debug(msg)
}

// Error ...
func Error(ctx context.Context, msg string, err error) {

	if err != nil {
		rqID := ctx.Value("request_id").(string)
		logger.Warn(
			msg+" :: "+err.Error(),
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

	logger.Warn(
		msg,
		zap.String("request_id", rqID),
		zap.String("hostname", hostName),
		zap.String("service", serviceName),
		zap.ByteString("request_body", contractReq),
		zap.ByteString("response_body", contractBody),
	)
}
