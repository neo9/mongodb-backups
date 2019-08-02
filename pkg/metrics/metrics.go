package metrics

import "github.com/prometheus/client_golang/prometheus"

type BackupMetrics struct {
	Total     *prometheus.CounterVec
	RetentionTotal     *prometheus.CounterVec
	Size      *prometheus.GaugeVec
	Duration  *prometheus.HistogramVec
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

	prom.Size = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "backup_size",
			Help:      "The size of backup.",
		},
		[]string{"name"},
	)

	prom.Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "backup_duration",
			Help:      "The duration of backup.",
		},
		[]string{"name"},
	)

	prometheus.MustRegister(prom.Total)
	prometheus.MustRegister(prom.RetentionTotal)
	prometheus.MustRegister(prom.Size)
	prometheus.MustRegister(prom.Duration)

	return prom
}