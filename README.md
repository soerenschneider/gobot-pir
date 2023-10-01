# gobot-pir
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-pir)](https://goreportcard.com/report/github.com/soerenschneider/gobot-pir)
![test-workflow](https://github.com/soerenschneider/gobot-pir/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/gobot-pir/actions/workflows/release.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/gobot-pir/actions/workflows/golangci-lint.yaml/badge.svg)

Detects and forwards motion events using a [PIR sensor](https://gobot.io/documentation/drivers/pir-motion-sensor/) and a Raspberry PI

## Features

🤖 Integrates with Home-Assistant<br/>
📊 Calculates statistics about motion events over time windows, accessible via MQTT and metrics<br/>
🔐 Allows connecting to secure MQTT brokers using TLS client certificates<br/>
🔭 Expose PIR events as metrics to enable alerting and Grafana dashboards<br/>

## Installation

### Binaries
Download a prebuilt binary from the [releases section](https://github.com/soerenschneider/gobot-pir/releases) for your system.

### From Source
As a prerequisite, you need to have [Golang SDK](https://go.dev/dl/) installed. Then you can install gobot-pir from source by invoking:
```shell
$ go install github.com/soerenschneider/gobot-pir@latest
```

## Configuration

gobot-pir can be fully configured using either environment variables or a config file.

| Field Name       | Type         | JSON Key       | Env Variable                  | Validation           | Default       | Description                      |
|------------------|--------------|----------------|-------------------------------|----------------------|---------------|----------------------------------|
| Placement        | string       | placement      | GOBOT_PIR_PLACEMENT           | required             | -             | Placement configuration          |
| MetricsAddr      | string       | metrics_addr   | GOBOT_PIR_METRICS_LISTEN_ADDR | omitempty,tcp_addr   | 0.0.0.0:9191  | Metrics address                  |
| IntervalSecs     | int          | interval_s     | GOBOT_PIR_INTERVAL_S          | min=1,max=300        |               | Interval in seconds              |
| StatIntervals    | []int        | stat_intervals | GOBOT_PIR_STAT_INTERVALS      | dive,min=10,max=3600 |               | Statistic intervals              |
| LogSensor        | bool         | log_sensor     | GOBOT_PIR_LOG_SENSOR_READINGS | false                |               | Log sensor readings              |
| MessageOn        | string       | message_on     | GOBOT_PIR_MSG_ON              |                      | ON            | Message when event is registered |
| MessageOff       | string       | message_off    | GOBOT_PIR_MSG_OFF             |                      | OFF           | Message when event stops         |
| MqttConfig       | MqttConfig   | -              |                               |                      |               | MQTT configuration               |
| SensorConfig     | SensorConfig |                |                               |                      |               | Sensor configuration             |


| Field Name       | Type     | JSON Key           | Env Variable                       | Validation                                       | Default | Description               |
|------------------|----------|--------------------|------------------------------------|--------------------------------------------------|---------|---------------------------|
| Disabled         | bool     | disable_mqtt       | GOBOT_PIR_MQTT_DISABLED            |                                                  | N/A     | Disabled MQTT             |
| Host             | string   | mqtt_host          | GOBOT_PIR_MQTT_BROKER              | required_if=Disabled false,mqtt_broker           | N/A     | MQTT broker host          |
| Topic            | string   | mqtt_topic         | GOBOT_PIR_MQTT_TOPIC               | required_if=Disabled false,mqtt_topic            | N/A     | MQTT topic                |
| StatsTopic       | string   | mqtt_stats_topic   | GOBOT_PIR_MQTT_STATS_TOPIC         | omitempty,mqtt_topic                             | N/A     | MQTT statistics topic     |
| ClientKeyFile    | string   | mqtt_ssl_key_file  | GOBOT_PIR_MQTT_TLS_CLIENT_KEY_FILE | required_unless=ClientCertFile '',omitempty,file | N/A     | MQTT client SSL key file  |
| ClientCertFile   | string   | mqtt_ssl_cert_file | GOBOT_PIR_MQTT_TLS_CLIENT_CRT_FILE | required_unless=ClientKeyFile '',omitempty,file  | N/A     | MQTT client SSL cert file |
| ServerCaFile     | string   | mqtt_ssl_ca_file   | GOBOT_PIR_MQTT_TLS_SERVER_CA_FILE  | omitempty,file                                   | N/A     | MQTT server SSL CA file   |


| Field Name               | Type     | JSON Key                 | Env Variable                       | Validation             | Default | Description                    |
|--------------------------|----------|--------------------------|------------------------------------|------------------------|---------|--------------------------------|
| GpioPin                  | string   | gpio_pin                 | GOBOT_PIR_GPIO_PIN                 | required,min=0         | 7       | GPIO pin configuration         |
| GpioPollingIntervalMs    | int      | gpio_polling_interval_ms | GOBOT_PIR_GPIO_POLLING_INTERVAL_MS | required,min=5,max=500 | 75      | GPIO polling interval in ms    |


### Via Config File

## Metrics

This project exposes the following metrics in Open Metrics format under the prefix `gobot_pir`

| Variable Name                      | Metric Type        | Description                                    | Labels              |
|------------------------------------|--------------------|------------------------------------------------|---------------------|
| version                            | GaugeVec           | Version information of this robot              | version, commit     |
| heartbeat_timestamp_seconds        | GaugeVec           | Heartbeat of this robot                        | placement           |
| motions_detected_total             | CounterVec         | Amount of motions detected                     | placement           |
| motions_detected_timestamp_seconds | GaugeVec           | Timestamp of latest motion detected            | placement           |
| messages_published_total           | CounterVec         | The assembleBot temperature in degrees Celsius | placement           |
| message_publish_errors_total       | CounterVec         | The assembleBot temperature in degrees Celsius | placement           |
| events_per_interval                | GaugeVec           | The number of events during given intervals    | interval, placement |
| slice_entries_total                | GaugeVec           | The amount of entries in the stats slice       | placement           |


## CHANGELOG
The changelog can be found [here](CHANGELOG.md)