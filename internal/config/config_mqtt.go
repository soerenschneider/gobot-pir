package config

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	// This regex is not a very strict check, we don't validate hostname or ip (v4, v6) addresses...
	mqttHostRegex = regexp.MustCompile(`^\w{3,}://.{3,}:\d{2,5}$`)

	// We don't care that technically it's allowed to start with a slash
	mqttTopicRegex = regexp.MustCompile(`^([\w%]+)(/[\w%]+)*$`)
)

type MqttConfig struct {
	Host           string `json:"mqtt_host,omitempty"`
	Topic          string `json:"mqtt_topic,omitempty"`
	ClientKeyFile  string `json:"mqtt_ssl_key_file,omitempty""`
	ClientCertFile string `json:"mqtt_ssl_cert_file,omitempty""`
	ServerCaFile   string `json:"mqtt_ssl_ca_file,omitempty"`
	StatsTopic     string `json:"mqtt_stats_topic,omitempty"`
}

func (conf *MqttConfig) Validate() error {
	if err := matchTopic(conf.Topic); err != nil {
		return errors.New("invalid mqtt topic provided")
	}

	if err := matchHost(conf.Host); err != nil {
		return err
	}

	return nil
}

func (conf *MqttConfig) UsesSslCerts() bool {
	return len(conf.ClientCertFile) > 0 && len(conf.ClientKeyFile) > 0
}

func (conf *MqttConfig) Print() {
	log.Printf("Host=%s", conf.Host)
	log.Printf("Topic=%s", conf.Topic)
	if len(conf.StatsTopic) > 0 {
		log.Printf("StatsTopic=%s", conf.Topic)
	}
	if conf.UsesSslCerts() {
		log.Printf("ClientCertificateFile=%s", conf.ClientCertFile)
	}
	if conf.UsesSslCerts() {
		log.Printf("ClientKeyFile=%s", conf.ClientKeyFile)
	}
	if len(conf.ServerCaFile) > 0 {
		log.Printf("ServerCaFile=%s", conf.ServerCaFile)
	}
}

func matchTopic(topic string) error {
	if !mqttTopicRegex.MatchString(topic) {
		return fmt.Errorf("invalid topic format used")
	}
	return nil
}

func matchHost(host string) error {
	if !mqttHostRegex.Match([]byte(host)) {
		return fmt.Errorf("invalid host format used")
	}
	return nil
}

func (conf *Config) FormatTopic() {
	if strings.Contains(conf.Topic, "%s") {
		conf.Topic = fmt.Sprintf(conf.Topic, conf.Placement)
	}
}
