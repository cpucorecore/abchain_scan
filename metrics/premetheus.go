package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	defaultMaxAge     = time.Second * 10
	defaultAgeBuckets = uint32(100)
	defaultObjectives = map[float64]float64{
		0.99: 0.01,
	}
)

var (
	CurrentHeight = prometheus.NewGauge(prometheus.GaugeOpts{Name: "current_height"})
	NewestHeight  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "newest_height"})
	TxCntByBlock  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "tx_cnt_by_block"})

	GetBlockDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_duration_ms",
		Help:       "get_block duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	GetBlockReceiptsDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_receipts_duration_ms",
		Help:       "get block receipts duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	BlockDelay = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "block_delay",
		Help:       "block delay in Seconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	BlockQueueSize = prometheus.NewGauge(prometheus.GaugeOpts{Name: "block_queue_size"})

	ParseBlockDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "parse_block_duration_ms",
		Help:       "parse block duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	DbOperationDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "db_operation_duration_ms",
		Help:       "db operation duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	SendBlockKafkaDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "send_block_kafka_duration_ms",
		Help:       "send block kafka duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	CallContractDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "call_contract_duration_ms",
		Help:       "call contract duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	CallContractArchiveDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "call_contract_archive_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	CallContractErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "call_contract_errors_total",
		},
		[]string{"is_retryable"},
	)

	GetPairDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_pair_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	GetTokenDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_token_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	VerifyPairDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "verify_pair_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	VerifyPairTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "verify_pair_total",
		},
		[]string{"result"},
	)

	VerifyPairOkByProtocol = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "verify_pair_ok_by_protocol_total",
		},
		[]string{"protocol"},
	)
)

func init() {
	prometheus.MustRegister(CurrentHeight)
	prometheus.MustRegister(NewestHeight)
	prometheus.MustRegister(TxCntByBlock)

	prometheus.MustRegister(GetBlockDurationMs)
	prometheus.MustRegister(GetBlockReceiptsDurationMs)
	prometheus.MustRegister(BlockDelay)
	prometheus.MustRegister(BlockQueueSize)

	prometheus.MustRegister(ParseBlockDurationMs)
	prometheus.MustRegister(DbOperationDurationMs)
	prometheus.MustRegister(SendBlockKafkaDurationMs)

	prometheus.MustRegister(CallContractDurationMs)
	prometheus.MustRegister(CallContractArchiveDurationMs)
	prometheus.MustRegister(CallContractErrors)
	prometheus.MustRegister(GetPairDurationMs)
	prometheus.MustRegister(GetTokenDurationMs)

	prometheus.MustRegister(VerifyPairDurationMs)
	prometheus.MustRegister(VerifyPairTotal)
	prometheus.MustRegister(VerifyPairOkByProtocol)
}

func init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf("%s:%d", "0.0.0.0", 9101), nil)
	}()
}
