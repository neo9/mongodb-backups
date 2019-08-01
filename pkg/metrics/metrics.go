package metrics

import "github.com/prometheus/client_golang/prometheus"

type BackupMetrics struct {
	Total   *prometheus.CounterVec
	Size    *prometheus.GaugeVec
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

	prom.Size = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "backup_size",
			Help:      "The size of backup.",
		},
		[]string{"name", "status"},
	)

	prometheus.MustRegister(prom.Total)
	prometheus.MustRegister(prom.Size)

	return prom
}