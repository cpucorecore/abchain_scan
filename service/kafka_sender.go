package service

import (
	"abchain_scan/config"
	"abchain_scan/log"
	"abchain_scan/metrics"
	"abchain_scan/types"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

type KafkaSender interface {
	Send(block *types.BlockInfo) error
}

type kafkaSender struct {
	ID            string
	conf          *config.KafkaConf
	sendTimeout   time.Duration
	asyncProducer sarama.AsyncProducer
}

func NewKafkaSender(conf *config.KafkaConf) KafkaSender {
	client := &kafkaSender{
		conf:        conf,
		sendTimeout: time.Millisecond * time.Duration(conf.SendTimeoutByMs),
	}

	sc := sarama.NewConfig()
	sc.Net.TLS.Enable = false
	sc.Producer.Return.Errors = true
	sc.Producer.RequiredAcks = sarama.WaitForLocal
	sc.Producer.Compression = sarama.CompressionSnappy
	sc.Producer.Flush.Frequency = 100 * time.Millisecond
	sc.Producer.Retry.Max = 10

	asyncProducer, err := sarama.NewAsyncProducer(conf.Brokers, sc)
	if err != nil {
		log.Logger.Fatal("kafka NewAsyncProducer err", zap.Error(err))
	}
	client.asyncProducer = asyncProducer
	client.processErrors()

	return client
}

func (s *kafkaSender) Close() {
	_ = s.asyncProducer.Close()
}

func (s *kafkaSender) processErrors() {
	errCh := s.asyncProducer.Errors()
	go func() {
		for {
			err, ok := <-errCh
			if !ok {
				log.Logger.Info("kafka asyncProducer error @ done", zap.Error(err))
				return
			}
			log.Logger.Info("kafka asyncProducer error", zap.Error(err))
		}
	}()
}

func (s *kafkaSender) Send(block *types.BlockInfo) error {
	if !s.conf.Enabled {
		return nil
	}

	data, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %v, %v", err, block)
	}

	now := time.Now()
	s.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: s.conf.Topic,
		Value: sarama.ByteEncoder(data),
	}
	metrics.SendBlockKafkaDurationMs.Observe(float64(time.Since(now).Milliseconds()))

	return nil
}
