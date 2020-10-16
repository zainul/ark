package nsq

import (
	"log"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
)

type Process interface {
	HandleMsg(msg *nsq.Message) error
}

// ConsumerConfig denotes all config available for  consumer.
type ConsumerConfig struct {
	MaxInFlight     int
	NumOfConsumers  int
	MaxAttempts     int
	RequeueInterval int
	Timeout         int
	Topic           string
	Channel         string
	LogPrefix       string
}

const (
	// Below are the default configs of conversion event consumer.
	DefaultMaxInflight     = 200
	DefaultMaxAttempts     = 3
	DefaultMsgTimeout      = 120
	DefaultRequeueInterval = 120
	DefaultNumOfConsumers  = 1
)

// OptionEvent sets options for option consumer.
type OptionEvent func(dgc *Consumer)

// Consumer implements handle message for nsq message.
type Consumer struct {
	cfg     ConsumerConfig
	address string
	process Process
}

// WithEventSettings sets Consumer settings.
func WithEventSettings(cfg ConsumerConfig) OptionEvent {
	return func(c *Consumer) {
		if cfg.MaxInFlight > 0 {
			c.cfg.MaxInFlight = cfg.MaxInFlight
		}
		if cfg.MaxAttempts > 0 {
			c.cfg.MaxAttempts = cfg.MaxAttempts
		}
		if cfg.Timeout > 0 {
			c.cfg.Timeout = cfg.Timeout
		}
		if cfg.RequeueInterval > 0 {
			c.cfg.RequeueInterval = cfg.RequeueInterval
		}
		if cfg.NumOfConsumers > 0 {
			c.cfg.NumOfConsumers = cfg.NumOfConsumers
		}
	}
}

//New initiate consumer of listen to topic
func New(address string, process Process, options ...OptionEvent) *Consumer {
	arl := &Consumer{
		cfg: ConsumerConfig{
			MaxInFlight:     DefaultMaxInflight,
			NumOfConsumers:  DefaultNumOfConsumers,
			MaxAttempts:     DefaultMaxAttempts,
			RequeueInterval: DefaultRequeueInterval,
			Timeout:         DefaultMsgTimeout,
		},
		address: address,
		process: process,
	}

	for _, opt := range options {
		opt(arl)
	}

	return arl
}

// getConfig returns config of AdScoreRecalculationLoggerConsumer
func (c *Consumer) getConfig() *nsq.Config {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = c.cfg.MaxInFlight
	cfg.MaxAttempts = uint16(c.cfg.MaxAttempts)
	cfg.MsgTimeout = time.Duration(c.cfg.Timeout) * time.Second
	cfg.DefaultRequeueDelay = time.Duration(c.cfg.RequeueInterval) * time.Second
	cfg.MaxBackoffDuration = 0
	return cfg
}

// Start starts ad score recalculation logger event consumer
func (c *Consumer) Start() error {
	q, err := nsq.NewConsumer(c.cfg.Topic, c.cfg.Channel, c.getConfig())
	if err != nil {
		return err
	}

	q.AddConcurrentHandlers(c, c.cfg.NumOfConsumers)
	q.SetLogger(log.New(os.Stderr, c.cfg.LogPrefix, log.Ltime), nsq.LogLevelError)

	return q.ConnectToNSQLookupd(c.address)
}

//HandleMessage handler nsq msg of Consumer
func (c *Consumer) HandleMessage(msg *nsq.Message) error {
	return c.process.HandleMsg(msg)
}
