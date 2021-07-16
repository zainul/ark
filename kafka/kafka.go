package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers      []string      `yaml:"connection"`
	GroupID      string        `yaml:"group-id"`
	Topic        string        `yaml:"topic"`
	URL          string        `yaml:"url"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
}

type Kafka struct {
	KW  *kafka.Writer
	KR  *kafka.Reader
	KCG *kafka.ConsumerGroup
}

func GetKafka(config Config) (Kafka, error) {
	kcg, err := getNewConsumerGroup(config.Brokers, config.GroupID, config.Topic)
	if err != nil {
		return Kafka{}, err
	}
	return Kafka{
		KW:  getWriterConfig(config.Brokers, config.Topic, config.WriteTimeout, config.ReadTimeout),
		KR:  getReaderConfig(config.Brokers, config.GroupID, config.Topic),
		KCG: kcg,
	}, nil
}

//err = bus.w.WriteMessages(context.Background(), kafka.Message{
//	Key: []byte(e.EventType), // 3 means the total number of brokers
//	// create an  message payload for the value
//	Value: msg,
//})
//if err != nil {
//	log.StdError(context.Background(), log.KV{}, err, "failed publish event")
//	return err
//}

func (kw Kafka) WriteMessage(ctx context.Context, key string, value string) error {
	return kw.KW.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})
}

//	for {
//		ctx := context.Background()
//		log.StdInfo(context.Background(), log.KV{"config reader ": c.r.Config().Brokers}, nil, "Start read kafka message")
//		msg, err := c.r.ReadMessage(ctx)
//		if err != nil {
//			log.StdError(ctx, log.KV{"topic": msg.Topic, "key": msg.Key}, err, "failed read message")
//			continue
//		}
//		var e event.Event
//		err = json.Unmarshal(msg.Value, &e)
//		if err != nil {
//			log.StdError(ctx, log.KV{"message": string(msg.Value)}, err, "failed parsing message")
//			continue
//		}
//		log.StdDebug(ctx, log.KV{"event": e}, nil, "message received")
//		handlers, ok := c.handlers[e.EventType]
//		if ok {
//			for handler := range handlers {
//				log.StdDebug(ctx, log.KV{"event": e}, nil, "message start handling")
//				handler.Handle(e)
//				log.StdDebug(ctx, log.KV{"event": e}, nil, "message done handling")
//			}
//		} else {
//			log.StdError(ctx, log.KV{"event": e.EventType}, errors.New("handler not found for event type"), "unhandled message")
//		}
//	}

func (kw Kafka) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return kw.KR.ReadMessage(ctx)
}

func getWriterConfig(kafkaBrokerUrls []string, topic string, writeTimeout, readTimeout time.Duration) (w *kafka.Writer) {
	return &kafka.Writer{
		Addr:         kafka.TCP(kafkaBrokerUrls...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
}

func getReaderConfig(brokers []string, groupID string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
}

func getNewConsumerGroup(brokers []string, groupID string, topic string) (*kafka.ConsumerGroup, error) {
	return kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      groupID,
		Brokers: brokers,
		Topics:  []string{topic},
	})
}
