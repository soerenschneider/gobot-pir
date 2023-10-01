package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/soerenschneider/gobot-pir/internal/config"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/mqtt"
)

const (
	heartbeatInterval = time.Duration(60) * time.Second
)

type PirSensor interface {
	gobot.Eventer
	Start() (err error)
	Halt() (err error)
	Name() string
	SetName(n string)
	Pin() string
	Connection() gobot.Connection
}

type MotionDetection struct {
	Driver      *gpio.PIRMotionDriver
	Adaptor     gobot.Connection
	MqttAdaptor *mqtt.Adaptor
	Stats       *SensorStats
	errorCount  int

	Config *config.Config
}

func (m *MotionDetection) publishMessage(msg []byte) {
	success := m.MqttAdaptor.Publish(m.Config.Topic, msg)
	if success {
		metricsMessagesPublished.WithLabelValues(m.Config.Placement).Inc()
	} else {
		metricsMessagePublishErrors.WithLabelValues(m.Config.Placement).Inc()
	}
}

func (m *MotionDetection) publishStatsMessage(stats map[string]int) error {
	msg, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("could not publish stats message: %v", err)
	}

	success := m.MqttAdaptor.Publish(m.Config.StatsTopic, msg)
	if success {
		metricsMessagesPublished.WithLabelValues(m.Config.Placement).Inc()
	} else {
		metricsMessagePublishErrors.WithLabelValues(m.Config.Placement).Inc()
	}
	return nil
}

func (m *MotionDetection) motionDetected(data interface{}) {
	metricsMotionsDetected.WithLabelValues(m.Config.Placement).Inc()
	metricsMotionTimestamp.WithLabelValues(m.Config.Placement).SetToCurrentTime()
	m.Stats.NewEvent()
	if len(m.Config.MessageOn) > 0 {
		m.publishMessage([]byte(m.Config.MessageOn))
	}
	if m.Config.LogSensor {
		log.Println("Detected motion")
	}
}

func (m *MotionDetection) motionStopped(data interface{}) {
	if len(m.Config.MessageOff) > 0 {
		m.publishMessage([]byte(m.Config.MessageOff))
	}
	if m.Config.LogSensor {
		log.Println("Motion stopped")
	}
}

func (m *MotionDetection) onError(data interface{}) {
	if m.errorCount > 10 {
		log.Fatalf("Too many errors reading from sensor, shutting down")
	}

	m.errorCount += 1
	log.Printf("GPIO error: %v", data)
}

func (m *MotionDetection) Send() {
	statsDict := map[string]int{}
	for _, stat := range m.Config.StatIntervals {
		count := m.Stats.GetEventCountNewerThan(time.Duration(stat) * time.Second)
		key := fmt.Sprintf("%ds", stat)
		statsDict[key] = count
		metricsStats.WithLabelValues(key, m.Config.Placement).Set(float64(count))
	}

	max, _ := m.Config.GetStatIntervalMax()
	m.Stats.PurgeEventsBefore(time.Now().Add(time.Duration(-max) * time.Second))
	metricsStatsSliceSize.WithLabelValues(m.Config.Placement).Set(float64(m.Stats.GetStatsSliceSize()))
	if err := m.publishStatsMessage(statsDict); err != nil {
		log.Printf("could not publish message: %v", err)
	}
}

func AssembleBot(motion *MotionDetection) *gobot.Robot {
	versionInfo.WithLabelValues(BuildVersion, CommitHash).Set(1)
	work := func() {
		if err := motion.Driver.On(gpio.MotionDetected, motion.motionDetected); err != nil {
			log.Printf("error for '%s' event: %v", gpio.MotionDetected, err)
		}

		if err := motion.Driver.On(gpio.MotionStopped, motion.motionStopped); err != nil {
			log.Printf("error for '%s' event: %v", gpio.MotionStopped, err)
		}

		if err := motion.Driver.On(gpio.Error, motion.onError); err != nil {
			log.Printf("error for '%s' event: %v", gpio.Error, err)
		}

		if len(motion.Config.MqttConfig.StatsTopic) > 0 && len(motion.Config.StatIntervals) > 0 {
			min, _ := motion.Config.GetStatIntervalMin()
			gobot.Every(time.Duration(min)*time.Second, motion.Send)
		}

		gobot.Every(heartbeatInterval, func() {
			metricsHeartbeat.WithLabelValues(motion.Config.Placement).SetToCurrentTime()
		})
	}

	adaptors := []gobot.Connection{motion.Adaptor}
	if motion.MqttAdaptor != nil {
		adaptors = append(adaptors, motion.MqttAdaptor)
	}

	return gobot.NewRobot(config.BotName,
		adaptors,
		[]gobot.Device{motion.Driver},
		work,
	)
}
