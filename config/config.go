package config

import (
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"os"
	"time"
)

type LogConf struct {
	Async                      bool `json:"async"`
	AsyncBufferSizeByByte      int  `json:"async_buffer_size_by_byte"`
	AsyncFlushIntervalBySecond int  `json:"async_flush_interval_by_second"`
}

type ChainConf struct {
	Endpoint        string `json:"endpoint"`
	EndpointArchive string `json:"endpoint_archive"`
	WsEndpoint      string `json:"ws_endpoint"`
}

type RedisConf struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BlockGetterConf struct {
	PoolSize         int       `json:"pool_size"`
	QueueSize        int       `json:"queue_size"`
	StartBlockNumber uint64    `json:"start_block_number"`
	Retry            RetryConf `json:"retry"`
}

type BlockHandlerConf struct {
	PoolSize        int `json:"pool_size"`
	ParseTxPoolSize int `json:"parse_tx_pool_size"`
	QueueSize       int `json:"queue_size"`
}

type RetryConf struct {
	Attempts  uint `json:"attempts"`
	DelayMs   int  `json:"delay_ms"`
	TimeoutMs int  `json:"timeout_ms"`
}

func (rc *RetryConf) GetRetryParams() *RetryParams {
	return &RetryParams{
		Attempts: retry.Attempts(rc.Attempts),
		Delay:    retry.Delay(time.Duration(rc.DelayMs) * time.Millisecond),
		Timeout:  time.Duration(rc.TimeoutMs) * time.Millisecond,
	}
}

type RetryParams struct {
	Attempts retry.Option  `json:"attempts"`
	Delay    retry.Option  `json:"delay"`
	Timeout  time.Duration `json:"timeout"`
}

type PriceServiceConf struct {
	PoolSize int `json:"pool_size"`
}

type KafkaConf struct {
	Enabled           bool     `json:"enabled"`
	Brokers           []string `json:"brokers"`
	Topic             string   `json:"topic"`
	SendTimeoutByMs   int      `json:"send_timeout_by_ms"`
	MaxRetry          int      `json:"max_retry"`
	RetryIntervalByMs int      `json:"retry_interval_by_ms"`
}

type ContractCallerConf struct {
	Retry *RetryConf `json:"retry"`
}

type DBDatasourceConf struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

func (c *DBDatasourceConf) GetPostgresDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		c.Host, c.Username, c.Password, c.DBName, c.Port)
}

type DBConf struct {
	Enabled      bool              `json:"enabled"`
	DBDatasource *DBDatasourceConf `json:"db_datasource"`
}

type Config struct {
	Log               *LogConf            `json:"log"`
	Chain             *ChainConf          `json:"chain"`
	Redis             *RedisConf          `json:"redis"`
	BlockGetter       *BlockGetterConf    `json:"block_getter"`
	BlockHandler      *BlockHandlerConf   `json:"block_handler"`
	EnableSequencer   bool                `json:"enable_sequencer"`
	PriceService      *PriceServiceConf   `json:"price_service"`
	Kafka             *KafkaConf          `json:"kafka"`
	ContractCaller    *ContractCallerConf `json:"contract_caller"`
	TxDatabase        *DBConf             `json:"tx_database"`
	TokenPairDatabase *DBConf             `json:"token_pair_database"`
}

var (
	defaultConfig = Config{
		Log: &LogConf{
			Async:                      false,
			AsyncBufferSizeByByte:      1000000,
			AsyncFlushIntervalBySecond: 1,
		},
		Chain: &ChainConf{
			Endpoint:        "https://base-rpc.publicnode.com",
			EndpointArchive: "https://base-rpc.publicnode.com",
			WsEndpoint:      "wss://base-rpc.publicnode.com",
		},
		Redis: &RedisConf{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
		},
		BlockGetter: &BlockGetterConf{
			PoolSize:         1,
			QueueSize:        1,
			StartBlockNumber: 48000000,
			Retry: RetryConf{
				Attempts:  10,
				DelayMs:   100,
				TimeoutMs: 5000,
			},
		},
		BlockHandler: &BlockHandlerConf{
			PoolSize:  1,
			QueueSize: 1,
		},
		EnableSequencer: true,
		PriceService: &PriceServiceConf{
			PoolSize: 1,
		},
		Kafka: &KafkaConf{
			Enabled:           false,
			Brokers:           []string{"localhost:9092"},
			Topic:             "block",
			SendTimeoutByMs:   5000,
			MaxRetry:          10,
			RetryIntervalByMs: 100,
		},
		ContractCaller: &ContractCallerConf{
			Retry: &RetryConf{
				Attempts:  10,
				DelayMs:   100,
				TimeoutMs: 3000,
			},
		},
		TxDatabase: &DBConf{
			Enabled: false,
			DBDatasource: &DBDatasourceConf{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				DBName:   "test",
			},
		},
		TokenPairDatabase: &DBConf{
			Enabled: false,
			DBDatasource: &DBDatasourceConf{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				DBName:   "test",
			},
		},
	}

	G = defaultConfig
)

func LoadConfigFile(configFilePath string) error {
	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&G); err != nil {
		return err
	}

	return nil
}
