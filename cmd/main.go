package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/soerenschneider/gobot-pir/internal"
	"github.com/soerenschneider/gobot-pir/internal/config"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/mqtt"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

const (
	cliConfFile = "config"
	cliVersion  = "version"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, cliConfFile, "", "File to read configuration from")
	version := flag.Bool(cliVersion, false, "Print version and exit")

	flag.Parse()

	if *version {
		fmt.Printf("%s (revision %s)", internal.BuildVersion, internal.CommitHash)
		os.Exit(0)
	}

	log.Printf("Started %s, version %s, commit %s", config.BotName, internal.BuildVersion, internal.CommitHash)
	conf := loadConfig(configFile)
	conf.Print()
	log.Println("Validating config...")
	err := conf.Validate()
	if err != nil {
		log.Fatalf("Could not validate config: %v", err)
	}
	conf.FormatTopic()
	run(conf)
}

func run(conf config.Config) {
	if conf.MetricConfig != "" {
		go internal.StartMetricsServer(conf.MetricConfig)
	}

	raspberry := raspi.NewAdaptor()
	driver := gpio.NewPIRMotionDriver(raspberry, conf.GpioPin, time.Millisecond*time.Duration(conf.GpioPollingIntervalMs))
	clientId := fmt.Sprintf("%s_%s", config.BotName, conf.Placement)
	mqttAdaptor := mqtt.NewAdaptor(conf.MqttConfig.Host, clientId)
	mqttAdaptor.SetAutoReconnect(true)
	mqttAdaptor.SetQoS(1)

	if conf.MqttConfig.UsesSslCerts() {
		log.Println("Setting TLS client cert and key...")
		mqttAdaptor.SetClientCert(conf.MqttConfig.ClientCertFile)
		mqttAdaptor.SetClientKey(conf.MqttConfig.ClientKeyFile)

		if len(conf.MqttConfig.ServerCaFile) > 0 {
			log.Println("Setting server CA...")
			mqttAdaptor.SetServerCert(conf.MqttConfig.ServerCaFile)
		}
	}

	adaptors := &internal.MotionDetection{
		Driver:      driver,
		Adaptor:     raspberry,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	bot := internal.AssembleBot(adaptors)
	err := bot.Start()
	if err != nil {
		log.Fatalf("could not start bot: %v", err)
	}
}

func loadConfig(configFile string) config.Config {
	if configFile == "" {
		log.Println("Building config from env vars")
		return config.ConfigFromEnv()
	}

	log.Printf("Reading config from file %s", configFile)
	conf, err := config.ReadJsonConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read config from %s: %v", configFile, err)
	}
	if nil == conf {
		log.Fatalf("Received empty config, should not happen")
	}
	return *conf
}
