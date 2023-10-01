package config

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
	GpioPin               string `json:"gpio_pin,omitempty" env:"GPIO_PIN" validate:"required,min=0"`
	GpioPollingIntervalMs int    `json:"gpio_polling_interval_ms,omitempty" env:"GPIO_POLLING_INTERVAL_MS" validate:"required,min=5,max=500"`
}
