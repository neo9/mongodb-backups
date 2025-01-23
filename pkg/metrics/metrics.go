package metrics

import (
	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
)

type BackupMetrics struct {
	Total                  *prometheus.CounterVec
	RetentionTotal         *prometheus.CounterVec
	BucketCount            *prometheus.GaugeVec
	Size                   *prometheus.GaugeVec
	Latency                *prometheus.SummaryVec
	LastSuccessfulSnapshot *prometheus.GaugeVec
}

func New(namespace string, subsystem string) *BackupMetrics {
	prom := &BackupMetrics{}

	prom.Total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "backup_total",
			Help:      "The total number of backups.",
		},
		[]string{"name", "status"},
	)

	prom.RetentionTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "retention_total",
			Help:      "The total number of retention removal.",
		},
		[]string{"name", "status"},
	)

	prom.BucketCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "bucket_snapshot_count",
			Help:      "The total number of snapshots stored in bucket.",
		},
		[]string{"name"},
	)

	prom.Size = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "snapshot_size",
			Help:      "The size of backup.",
		},
		[]string{"name"},
	)

	prom.Latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "snapshot_latency",
			Help:      "Backup duration in seconds",
		},
		[]string{"name"},
	)

	prom.LastSuccessfulSnapshot = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "last_successful_snapshot",
			Help:      "The size of backup.",
		},
		[]string{"name"},
	)

	safeRegister(prom.Total)
	safeRegister(prom.RetentionTotal)
	safeRegister(prom.BucketCount)
	safeRegister(prom.Size)
	safeRegister(prom.Latency)
	safeRegister(prom.LastSuccessfulSnapshot)

	return prom
}

func safeRegister(collector prometheus.Collector) {
	err := prometheus.Register(collector)
	if err != nil {
		log.Warn("Metric %s already present in the system")
	}
}
