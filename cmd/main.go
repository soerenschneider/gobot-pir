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
	conf, err := config.Read(configFile)
	if err != nil {
		log.Fatalf("Could not read config: %v", err)
	}
	config.PrintFields(conf)
	log.Println("Validating config...")
	if err := config.Validate(conf); err != nil {
		log.Fatalf("Could not validate config: %v", err)
	}
	conf.FormatTopic()
	run(conf)
}

func run(conf *config.Config) {
	if conf.MetricsAddr != "" {
		go internal.StartMetricsServer(conf.MetricsAddr)
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
		Stats:       internal.NewSensorStats(),
	}

	bot := internal.AssembleBot(adaptors)
	err := bot.Start()
	if err != nil {
		log.Fatalf("could not start bot: %v", err)
	}
}
