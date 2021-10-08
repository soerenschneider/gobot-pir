package internal

import (
	"gobot-pir/internal/config"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/mqtt"
	"log"
	"time"
)

const (
	messageOnDetected = "ON"
	messageOnStopped  = "OFF"
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
		metricsMessagesPublished.WithLabelValues(m.Config.Location).Inc()
	} else {
		metricsMessagePublishErrors.WithLabelValues(m.Config.Location).Inc()
	}
}

func AssembleBot(motion *MotionDetection) *gobot.Robot {
	errorCnt := 0
	work := func() {
		gobot.Every(heartbeatInterval, func() {
			metricsHeartbeat.WithLabelValues(motion.Config.Location).SetToCurrentTime()
		})

		motion.Driver.On(gpio.MotionDetected, func(data interface{}) {
			metricsMotionsDetected.WithLabelValues(motion.Config.Location).Inc()
			metricsMotionTimestamp.WithLabelValues(motion.Config.Location).SetToCurrentTime()
			motion.publishMessage([]byte(messageOnDetected))
			if motion.Config.LogSensor {
				log.Println("Detected motion")
			}
		})

		motion.Driver.On(gpio.MotionStopped, func(data interface{}) {
			motion.publishMessage([]byte(messageOnStopped))
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
