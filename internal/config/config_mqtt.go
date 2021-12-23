package config

import (
	"errors"
	"log"
	"regexp"
)

var (
	// This regex is not a very strict check, we don't validate hostname or ip (v4, v6) addresses...
	mqttHostRegex = regexp.MustCompile(`^\w{3,}://.{3,}:\d{2,5}$`)

	// We don't care that technically it's allowed to start with a slash
	mqttTopicRegex = regexp.MustCompile("^([\\w%]+)(/[\\w%]+)*$")
)

type MqttConfig struct {
	Host       string `json:"mqtt_host,omitempty"`
	Topic      string `json:"mqtt_topic,omitempty"`
	StatsTopic string `json:"mqtt_stats_topic,omitempty"`
	Username   string `json:"mqtt_username,omitempty"`
	Password   string `json:"mqtt_password,omitempty"`
}

func (conf *MqttConfig) ConfigFromEnv() {
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

	mqttUsername, err := fromEnv("MQTT_USERNAME")
	if err == nil {
		conf.Username = mqttUsername
	}

	mqttPassword, err := fromEnv("MQTT_PASSWORD")
	if err == nil {
		conf.Password = mqttPassword
	}
}

func (conf *MqttConfig) Validate() error {
	if err := matchTopic(conf.Topic); err != nil {
		return errors.New("invalid mqtt topic provided")
	}

	if err := matchHost(conf.Host); err != nil {
		return err
	}

	hasAuthUsername := len(conf.Username) > 0
	hasAuthPassword := len(conf.Password) > 0
	if hasAuthUsername && !hasAuthPassword || !hasAuthUsername && hasAuthPassword {
		return errors.New("must either specify both username and password or no authentication options at all")
	}

	return nil
}

func (conf *MqttConfig) Print() {
	log.Printf("Host=%s", conf.Host)
	log.Printf("Topic=%s", conf.Topic)
	if conf.UsesAuth() {
		log.Printf("Username=%s", conf.Username)
		log.Println("Topic=*** (Redacted)")
	}
	if len(conf.StatsTopic) > 0 {
		log.Printf("StatsTopic=%s", conf.Topic)
	}
}

func (conf *MqttConfig) UsesAuth() bool {
	return len(conf.Username) > 0 && len(conf.Password) > 0
}
