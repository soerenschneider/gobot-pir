package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const BotName = "gobot_motion_detection"

// This regex is not a very strict check, we don't validate hostname or ip (v4, v6) addresses...
var mqttHostRegex = regexp.MustCompile(`\w{3,}://.{3,}:\d{2,4}`)

type Config struct {
	Location        string `json:"location,omitempty"`
	MetricConfig    string `json:"metrics_addr,omitempty"`
	Pin             string `json:"gpio_pin,omitempty"`
	PollingInterval int    `json:"polling_interval_ms,omitempty"`
	LogMotions      bool   `json:"log_motions,omitempty"`
	MqttConfig
}

type MqttConfig struct {
	Host     string `json:"mqtt_host,omitempty"`
	ClientId string `json:"mqtt_client_id,omitempty"`
	Topic    string `json:"mqtt_topic,omitempty"`
}

func DefaultConfig() Config {
	location := fromEnv(fmt.Sprintf("%s_LOCATION", strings.ToUpper(BotName)), "")
	return Config{
		Location:        location,
		LogMotions:      fromEnvBool(fmt.Sprintf("%s_LOG_MOTIONS", strings.ToUpper(BotName)), false),
		Pin:             fromEnv(fmt.Sprintf("%s_GPIO_PIN", strings.ToUpper(BotName)), "7"),
		PollingInterval: fromEnvInt(fmt.Sprintf("%s_GPIO_POLLING_MS", strings.ToUpper(BotName)), 75),
		MqttConfig: MqttConfig{
			Host:     fromEnv(fmt.Sprintf("%s_MQTT_HOST", strings.ToUpper(BotName)), ""),
			ClientId: fromEnv(fmt.Sprintf("%s_MQTT_CLIENT_ID", strings.ToUpper(BotName)), fmt.Sprintf("%s-%s", BotName, location)),
			Topic:    fromEnv(fmt.Sprintf("%s_MQTT_TOPIC", strings.ToUpper(BotName)), fmt.Sprintf("%s/%s", BotName, location)),
		},
		MetricConfig: fromEnv(fmt.Sprintf("%s_METRICS_ADDR", strings.ToUpper(BotName)), ":9191"),
	}
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

func (c *Config) Validate() error {
	if len(c.Location) == 0 {
		return fmt.Errorf("empty location provided")
	}

	parsedPin, err := strconv.Atoi(c.Pin)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as pin: %v", c.Pin, err)
	}
	if parsedPin < 0 {
		return fmt.Errorf("invalid pin provided: %d", parsedPin)
	}

	if c.PollingInterval < 5 {
		return fmt.Errorf("polling interval can not be smaller than 5: %d", c.PollingInterval)
	}

	if c.PollingInterval > 500 {
		return fmt.Errorf("polling interval too high: %d", c.PollingInterval)
	}

	// TODO: improve check
	if strings.Index(c.MqttConfig.Topic, " ") != -1 {
		return fmt.Errorf("invalid mqtt topic provided")
	}

	return matchHost(c.MqttConfig.Host)
}

func (c *Config) Print() {
	log.Printf("Location=%s", c.Location)
	log.Printf("LogMotions=%t", c.LogMotions)
	log.Printf("MetricConfig=%s", c.MetricConfig)
	log.Printf("GpioPin=%s", c.Pin)
	log.Printf("GpioPollingIntervalMs=%d", c.PollingInterval)
	log.Printf("Host=%s", c.Host)
	log.Printf("Topic=%s", c.Topic)
	log.Printf("ClientId=%s", c.ClientId)
}

func matchHost(host string) error {
	if !mqttHostRegex.Match([]byte(host)) {
		return fmt.Errorf("invalid host format used")
	}
	return nil
}

func fromEnv(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

func fromEnvInt(name string, def int) int {
	val := os.Getenv(name)
	if val == "" {
		return def
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return parsed
}

func fromEnvBool(name string, def bool) bool {
	val := os.Getenv(name)
	if val == "" {
		return def
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return parsed
}
