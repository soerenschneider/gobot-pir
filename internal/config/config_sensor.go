package config

import (
	"fmt"
	"log"
	"strconv"
)

const (
	defaultGpioPin               = "7"
	defaultGpioPollingIntervalMs = 75
)

func defaultSensorConfig() SensorConfig {
	return SensorConfig{
		GpioPin:               defaultGpioPin,
		GpioPollingIntervalMs: defaultGpioPollingIntervalMs,
	}
}

type SensorConfig struct {
	GpioPin               string `json:"gpio_pin,omitempty"`
	GpioPollingIntervalMs int    `json:"gpio_polling_interval_ms,omitempty"`
}

func (conf *SensorConfig) Validate() error {
	parsedPin, err := strconv.Atoi(conf.GpioPin)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as pin: %v", conf.GpioPin, err)
	}
	if parsedPin < 0 {
		return fmt.Errorf("invalid pin provided: %d", parsedPin)
	}

	if conf.GpioPollingIntervalMs < 5 {
		return fmt.Errorf("polling interval must not be smaller than 5: %d", conf.GpioPollingIntervalMs)
	}

	if conf.GpioPollingIntervalMs > 500 {
		return fmt.Errorf("polling interval too high: %d", conf.GpioPollingIntervalMs)
	}

	return nil
}

func (conf *SensorConfig) ConfigFromEnv() {
	gpioPin, err := fromEnv("GPIO_PIN")
	if err == nil {
		conf.GpioPin = gpioPin
	}

	gpioPollingInterval, err := fromEnvInt("GPIO_POLLING_INTERVAL_MS")
	if err == nil {
		conf.GpioPollingIntervalMs = gpioPollingInterval
	}
}

func (conf *SensorConfig) Print() {
	log.Printf("GpioPin=%s", conf.GpioPin)
	log.Printf("GpioPollingIntervalMs=%d", conf.GpioPollingIntervalMs)
}
