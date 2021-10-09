package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soerenschneider/gobot-pir/internal/config"
	"log"
	"net/http"
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
	}, []string{"location"})

	metricsMotionsDetected = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "motions_detected_total",
		Subsystem: "sensor",
		Help:      "Amount of motions detected",
	}, []string{"location"})

	metricsMotionTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "motions_detected_timestamp_seconds",
		Subsystem: "sensor",
		Help:      "Timestamp of latest motion detected",
	}, []string{"location"})

	metricsMessagesPublished = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "messages_published_total",
		Subsystem: "mqtt",
		Help:      "The assembleBot temperature in degrees Celsius",
	}, []string{"location"})

	metricsMessagePublishErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "message_publish_errors_total",
		Subsystem: "mqtt",
		Help:      "The assembleBot temperature in degrees Celsius",
	}, []string{"location"})
)

func StartMetricsServer(listenAddr string) {
	log.Printf("Starting metrics listener at %s", listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Could not start metrics listener: %v", err)
	}
}
