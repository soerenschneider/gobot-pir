package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
)

const (
	BotName                = "gobot_pir"
	defaultLogValues       = false
	defaultIntervalSeconds = 30
	defaultMetricConfig    = "0.0.0.0:9191"
	defaultMessageOn       = "ON"
	defaultMessageOff      = "OFF"
	maxStatsBucketSeconds  = 7200
)

var (
	once                       sync.Once
	validate                   *validator.Validate
	defaultStatsBucketsSeconds = []int{5, 15, 30, 60, 120, 300, 600, 1800}
)

type Config struct {
	Placement     string `json:"placement,omitempty" env:"PLACEMENT" validate:"required"`
	MetricsAddr   string `json:"metrics_addr,omitempty" env:"METRICS_LISTEN_ADDR" validate:"omitempty,tcp_addr"`
	IntervalSecs  int    `json:"interval_s,omitempty" env:"INTERVAL_S" validate:"min=1,max=300"`
	StatIntervals []int  `json:"stat_intervals,omitempty" env:"STAT_INTERVALS" validate:"dive,min=10,max=3600"`
	LogSensor     bool   `json:"log_sensor,omitempty" env:"LOG_SENSOR_READINGS"`
	MessageOn     string `json:"message_on" env:"MSG_ON"`
	MessageOff    string `json:"message_off" env:"MSG_OFF"`
	MqttConfig
	SensorConfig
}

func DefaultConfig() Config {
	statInterval := make([]int, len(defaultStatsBucketsSeconds))
	copy(statInterval, defaultStatsBucketsSeconds)

	return Config{
		LogSensor:     defaultLogValues,
		IntervalSecs:  defaultIntervalSeconds,
		StatIntervals: statInterval,
		MetricsAddr:   defaultMetricConfig,
		SensorConfig:  defaultSensorConfig(),
		MessageOn:     defaultMessageOn,
		MessageOff:    defaultMessageOff,
	}
}

func Read(filePath string) (*Config, error) {
	ret := DefaultConfig()

	if len(filePath) > 0 {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("could not read config from file: %v", err)
		}

		err = json.Unmarshal(fileContent, &ret)
		if err != nil {
			return nil, err
		}
	}

	opts := env.Options{
		Prefix: fmt.Sprintf("%s_", strings.ToUpper(BotName)),
	}
	err := env.ParseWithOptions(&ret, opts)
	return &ret, err
}

func Validate(s interface{}) error {
	once.Do(func() {
		validate = validator.New()
		if err := validate.RegisterValidation("mqtt_topic", validateTopic); err != nil {
			log.Fatal("could not build custom validation 'mqtt_topic'")
		}
		if err := validate.RegisterValidation("mqtt_broker", validateBroker); err != nil {
			log.Fatal("could not build custom validation 'validateBroker'")
		}
	})
	return validate.Struct(s)
}

func validateTopic(fl validator.FieldLevel) bool {
	// Get the field value and check if it's a slice
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	topic, ok := field.Interface().(string)
	if !ok || !matchTopic(topic) {
		return false
	}

	return true
}

func validateBroker(fl validator.FieldLevel) bool {
	// Get the field value and check if it's a slice
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	// Convert to string and check its value
	broker, ok := field.Interface().(string)
	if !ok || !matchHost(broker) {
		return false
	}

	return true
}

func (conf *Config) GetStatIntervalMin() (int, error) {
	if len(conf.StatIntervals) == 0 {
		return -1, fmt.Errorf("empty array provided")
	}

	min := conf.StatIntervals[0]
	for _, val := range conf.StatIntervals {
		if val < min {
			min = val
		}
	}

	return min, nil
}

func (conf *Config) GetStatIntervalMax() (int, error) {
	if len(conf.StatIntervals) == 0 {
		return -1, fmt.Errorf("empty array provided")
	}

	max := conf.StatIntervals[0]
	for _, val := range conf.StatIntervals {
		if val > max {
			max = val
		}
	}

	return max, nil
}
