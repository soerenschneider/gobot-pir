package main

import (
	"flag"
	"fmt"
	"gobot-pir/internal"
	"gobot-pir/internal/config"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

func main() {
	log.Printf("Started %s, version %s, commit %s, built at %s", config.BotName, internal.BuildVersion, internal.CommitHash, internal.BuildTime)
	conf := getConfig()
	err := conf.Validate()
	conf.Print()
	if err != nil {
		log.Fatalf("Could not build config: %v", err)
	}

	if conf.MetricConfig != "" {
		go internal.StartMetricsServer(conf.MetricConfig)
	}

	raspberry := raspi.NewAdaptor()
	driver := gpio.NewPIRMotionDriver(raspberry, conf.GpioPin, time.Millisecond*time.Duration(conf.GpioPollingIntervalMs))
	clientId := fmt.Sprintf("%s_%s", config.BotName, conf.Location)
	mqttAdaptor := mqtt.NewAdaptor(conf.MqttConfig.Host, clientId)
	adaptors := &internal.MotionDetection{
		Driver:      driver,
		Adaptor:     raspberry,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	bot := internal.AssembleBot(adaptors)
	err = bot.Start()
	if err != nil {
		log.Fatalf("could not start bot: %v", err)
	}
}

func getConfig() config.Config {
	var configFile string
	flag.StringVar(&configFile, "config", "", "File to read configuration from")
	flag.Parse()
	if configFile == "" {
		log.Println("Building config from env vars")
		return config.DefaultConfig()
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
