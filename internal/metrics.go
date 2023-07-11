package internal

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soerenschneider/gobot-pir/internal/config"
)

const namespace = config.BotName

var (
	versionInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "version",
		Help:      "Version information of this robot",
	}, []string{"version", "commit"})

	metricsHeartbeat = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "heartbeat_timestamp_seconds",
		Help:      "Heartbeat of this robot",
	}, []string{"placement"})

	metricsMotionsDetected = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "motions_detected_total",
		Subsystem: "sensor",
		Help:      "Amount of motions detected",
	}, []string{"placement"})

	metricsMotionTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "motions_detected_timestamp_seconds",
		Subsystem: "sensor",
		Help:      "Timestamp of latest motion detected",
	}, []string{"placement"})

	metricsMessagesPublished = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "messages_published_total",
		Subsystem: "mqtt",
		Help:      "The assembleBot temperature in degrees Celsius",
	}, []string{"placement"})

	metricsMessagePublishErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "message_publish_errors_total",
		Subsystem: "mqtt",
		Help:      "The assembleBot temperature in degrees Celsius",
	}, []string{"placement"})

	metricsStats = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "events_per_interval",
		Subsystem: "stats",
		Help:      "The number of events during given intervals",
	}, []string{"interval", "placement"})

	metricsStatsSliceSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "slice_entries_total",
		Subsystem: "stats",
		Help:      "The amount of entries in the stats slice",
	}, []string{"placement"})
)

func StartMetricsServer(listenAddr string) {
	log.Printf("Starting metrics listener at %s", listenAddr)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not start metrics listener: %v", err)
	}
}
