package metrics

import "github.com/prometheus/client_golang/prometheus"

type BackupMetrics struct {
	Total     *prometheus.CounterVec
	RetentionTotal     *prometheus.CounterVec
	BucketCount  *prometheus.GaugeVec
	Size      *prometheus.GaugeVec
	Latency  *prometheus.SummaryVec
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
		[]string{"name", "status"},
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

	prometheus.MustRegister(prom.Total)
	prometheus.MustRegister(prom.RetentionTotal)
	prometheus.MustRegister(prom.BucketCount)
	prometheus.MustRegister(prom.Size)
	prometheus.MustRegister(prom.Latency)

	return prom
}