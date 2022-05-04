package internal

import (
	"encoding/json"
	"fmt"
	"github.com/soerenschneider/gobot-pir/internal/config"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/mqtt"
	"log"
	"time"
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

	Config config.Config
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

func AssembleBot(motion *MotionDetection) *gobot.Robot {
	versionInfo.WithLabelValues(BuildVersion, CommitHash).Set(1)
	errorCnt := 0
	stats := NewSensorStats()
	work := func() {
		motion.Driver.On(gpio.MotionDetected, func(data interface{}) {
			metricsMotionsDetected.WithLabelValues(motion.Config.Placement).Inc()
			metricsMotionTimestamp.WithLabelValues(motion.Config.Placement).SetToCurrentTime()
			stats.NewEvent()
			if len(motion.Config.MessageOn) > 0 {
				motion.publishMessage([]byte(motion.Config.MessageOn))
			}
			if motion.Config.LogSensor {
				log.Println("Detected motion")
			}
		})

		motion.Driver.On(gpio.MotionStopped, func(data interface{}) {
			if len(motion.Config.MessageOff) > 0 {
				motion.publishMessage([]byte(motion.Config.MessageOff))
			}
			if motion.Config.LogSensor {
				log.Println("Motion stopped")
			}
		})

		motion.Driver.On(gpio.Error, func(data interface{}) {
			if errorCnt > 10 {
				log.Fatalf("Too many errors reading from sensor, shutting down")
			}
			errorCnt += 1
			log.Printf("GPIO error: %v", data)
		})

		gobot.Every(heartbeatInterval, func() {
			metricsHeartbeat.WithLabelValues(motion.Config.Placement).SetToCurrentTime()
		})

		if len(motion.Config.MqttConfig.StatsTopic) != 0 && len(motion.Config.StatIntervals) > 0 {
			min, _ := motion.Config.GetStatIntervalMin()
			max, _ := motion.Config.GetStatIntervalMax()

			gobot.Every(time.Duration(min)*time.Second, func() {
				statsDict := map[string]int{}
				for _, stat := range motion.Config.StatIntervals {
					count := stats.GetEventCountNewerThan(time.Duration(stat) * time.Second)
					key := fmt.Sprintf("%ds", stat)
					statsDict[key] = count
					metricsStats.WithLabelValues(key, motion.Config.Placement).Set(float64(count))
				}

				stats.PurgeEventsBefore(time.Now().Add(time.Duration(-max) * time.Second))
				metricsStatsSliceSize.WithLabelValues(motion.Config.Placement).Set(float64(stats.GetStatsSliceSize()))
				motion.publishStatsMessage(statsDict)
			})
		}
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
