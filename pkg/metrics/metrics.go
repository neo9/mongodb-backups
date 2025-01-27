package metrics

import (
	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
)

type BackupMetrics struct {
	Total                  *prometheus.CounterVec
	BackupTries            *prometheus.CounterVec
	RetentionTotal         *prometheus.CounterVec
	RetentionBucketCount   *prometheus.GaugeVec
	BackupSize             *prometheus.GaugeVec
	SnapshotLatency        *prometheus.SummaryVec
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
			Help:      "The total number of retention removal tries.",
		},
		[]string{"name", "status"},
	)

	prom.RetentionBucketCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "bucket_snapshot_count",
			Help:      "The total number of snapshots stored in bucket.",
		},
		[]string{"name"},
	)

	prom.BackupSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "snapshot_size",
			Help:      "The size of backup.",
		},
		[]string{"name"},
	)

	prom.SnapshotLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "snapshot_latency",
			Help:      "The latency to create a backup in seconds.",
		},
		[]string{"name"},
	)

	prom.SnapshotLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "snapshot_latency",
			Help:      "The latency to create a backup in seconds.",
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

	prom.BackupTries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "number of backup tries",
			Help:      "The number of tries for each backup.",
		},
		[]string{"name", "status"},
	)

	safeRegister("Total", prom.Total)
	safeRegister("RetentionTotal", prom.RetentionTotal)
	safeRegister("BucketCount", prom.RetentionBucketCount)
	safeRegister("Size", prom.BackupSize)
	safeRegister("Latency", prom.SnapshotLatency)
	safeRegister("LastSuccessfulSnapshot", prom.LastSuccessfulSnapshot)
	safeRegister("BackupTries", prom.BackupTries)

	return prom
}

func safeRegister(metric_name string, collector prometheus.Collector) {
	err := prometheus.Register(collector)
	if err != nil {
		log.Warn("Metric %s already present in the system", metric_name)
	}
}
