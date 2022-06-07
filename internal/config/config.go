package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	BotName                = "gobot_pir"
	defaultLogValues       = false
	defaultIntervalSeconds = 30
	defaultMetricConfig    = ":9191"
	defaultMessageOn       = "ON"
	maxStatsBucketSeconds  = 7200
)

var (
	defaultStatsBucketsSeconds = []int{5, 15, 30, 60, 120, 300, 600, 1800}
)

type Config struct {
	Placement     string `json:"placement,omitempty"`
	MetricConfig  string `json:"metrics_addr,omitempty"`
	IntervalSecs  int    `json:"interval_s,omitempty"`
	StatIntervals []int  `json:"stat_intervals,omitempty"`
	LogSensor     bool   `json:"log_sensor,omitempty"`
	MessageOn     string `json:"message_on"`
	MessageOff    string `json:"message_off"`
	MqttConfig
	SensorConfig
}

func DefaultConfig() Config {
	return Config{
		LogSensor:     defaultLogValues,
		IntervalSecs:  defaultIntervalSeconds,
		StatIntervals: defaultStatsBucketsSeconds,
		MetricConfig:  defaultMetricConfig,
		SensorConfig:  defaultSensorConfig(),
		MessageOn:     "ON",
	}
}

func ConfigFromEnv() Config {
	conf := DefaultConfig()

	placement, err := fromEnv("PLACEMENT")
	if err == nil {
		conf.Placement = placement
	}

	logSensor, err := fromEnvBool("LOG_SENSOR")
	if err == nil {
		conf.LogSensor = logSensor
	}

	intervalSeconds, err := fromEnvInt("INTERVAL_S")
	if err == nil {
		conf.IntervalSecs = intervalSeconds
	}

	mqttHost, err := fromEnv("MQTT_HOST")
	if err == nil {
		conf.Host = mqttHost
	}

	mqttTopic, err := fromEnv("MQTT_TOPIC")
	if err == nil {
		conf.Topic = mqttTopic
	}

	mqttStatsTopic, err := fromEnv("MQTT_STATS_TOPIC")
	if err == nil {
		conf.StatsTopic = mqttStatsTopic
	}

	metricConfig, err := fromEnv("METRICS_ADDR")
	if err == nil {
		conf.MetricConfig = metricConfig
	}

	clientKeyFile, err := fromEnv("SSL_CLIENT_KEY_FILE")
	if err == nil {
		conf.ClientKeyFile = clientKeyFile
	}

	clientCertFile, err := fromEnv("SSL_CLIENT_CERT_FILE")
	if err == nil {
		conf.ClientCertFile = clientCertFile
	}

	conf.SensorConfig.ConfigFromEnv()

	return conf
}

func ReadJsonConfig(filePath string) (*Config, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config from file: %v", err)
	}

	ret := DefaultConfig()
	err = json.Unmarshal(fileContent, &ret)
	return &ret, err
}

func (conf *Config) Validate() error {
	if conf.Placement == "" {
		return errors.New("empty placement provided")
	}

	if conf.IntervalSecs < 30 {
		return fmt.Errorf("invalid interval: must not be lower than 30 but is %d", conf.IntervalSecs)
	}

	if conf.IntervalSecs > 300 {
		return fmt.Errorf("invalid interval: mut not be greater than 300 but is %d", conf.IntervalSecs)
	}

	if len(conf.StatIntervals) > 0 {
		min, _ := conf.GetStatIntervalMin()
		if min < 1 {
			return fmt.Errorf("minimal value in stats bucket must not be < 1: %d", min)
		}

		max, _ := conf.GetStatIntervalMax()
		if max > maxStatsBucketSeconds {
			return fmt.Errorf("maximal value in stats bucket must not be > %d: %d", maxStatsBucketSeconds, max)
		}
	}

	if err := conf.SensorConfig.Validate(); err != nil {
		return err
	}

	return conf.MqttConfig.Validate()
}

func (conf *Config) Print() {
	log.Println("-----------------")
	log.Println("Configuration:")
	log.Printf("Placement=%s", conf.Placement)
	log.Printf("LogSensor=%t", conf.LogSensor)
	log.Printf("MetricConfig=%s", conf.MetricConfig)
	log.Printf("IntervalSecs=%d", conf.IntervalSecs)
	log.Printf("Host=%s", conf.Host)
	log.Printf("Topic=%s", conf.Topic)

	if len(conf.StatIntervals) > 0 {
		log.Printf("StatIntervals=%v", conf.StatIntervals)
	}

	conf.SensorConfig.Print()
	conf.MqttConfig.Print()

	log.Println("-----------------")
}

func computeEnvName(name string) string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(BotName), strings.ToUpper(name))
}

func fromEnv(name string) (string, error) {
	name = computeEnvName(name)
	val := os.Getenv(name)
	if val == "" {
		return "", errors.New("not defined")
	}
	return val, nil
}

func fromEnvInt(name string) (int, error) {
	val, err := fromEnv(name)
	if err != nil {
		return -1, err
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return -1, err
	}
	return parsed, nil
}

func fromEnvBool(name string) (bool, error) {
	val, err := fromEnv(name)
	if err != nil {
		return false, err
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return parsed, nil
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
