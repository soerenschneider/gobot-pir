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

## MQTT Payloads

### Motion detected event
"ON"

### Motion stopped event
"OFF"

## Configuration

gobot-pir can be fully configured using either environment variables or a config file.

### Environment Variables Reference
| ENV                                     | Default                           | Description                                      |
|-----------------------------------------|-----------------------------------|--------------------------------------------------|
| GOBOT_MOTION_DETECTION_PLACEMENT        | -                                 | Location short name of this motion detection bot |
| GOBOT_MOTION_DETECTION_LOG_MOTIONS      | false                             | Write a log message for every motion event       |
| GOBOT_MOTION_DETECTION_GPIO_PIN         | 7                                 | GPIO pin to poll                                 |
| GOBOT_MOTION_DETECTION_GPIO_POLLING_MS  | 75                                | GPIO polling frequency in milliseconds           |
| GOBOT_MOTION_DETECTION_MQTT_HOST        | gobot_motion_detection-$PLACEMENT | MQTT connection broker                           |
| GOBOT_MOTION_DETECTION_MQTT_CLIENT_ID   | gobot_motion_detection-$PLACEMENT | Client ID for the MQTT connection                |
| GOBOT_MOTION_DETECTION_MQTT_TOPIC       | gobot_motion_detection/$PLACEMENT | Topic to publish messages into                   |
| GOBOT_MOTION_DETECTION_METRICS_ADDR     | :9400                             | Prometheus http handler listen address           |

### Via Config File

| Struct          | Field                     | Type          | JSON Tag                             | Optional | Defaults |
|-----------------|---------------------------|---------------|--------------------------------------|----------|----------|
| Config          | Placement                 | string        | "placement,omitempty"                | Yes      |          |
|                 | MetricsAddr               | string        | "metrics_addr,omitempty"             | Yes      | ":9191"  |
|                 | IntervalSecs              | int           | "interval_s,omitempty"               | Yes      | 30       |
|                 | StatIntervals             | []int         | "stat_intervals,omitempty"           | Yes      |          |
|                 | LogSensor                 | bool          | "log_sensor,omitempty"               | Yes      | false    |
|                 | MessageOn                 | string        | "message_on"                         | No       | "ON"     |
|                 | MessageOff                | string        | "message_off"                        | No       |          |
|                 | MqttConfig                | MqttConfig    |                                      | No       |          |
|                 | SensorConfig              | SensorConfig  |                                      | No       |          |
| MqttConfig      | Host                      | string        | "mqtt_host,omitempty"                | Yes      |          |
|                 | Topic                     | string        | "mqtt_topic,omitempty"               | Yes      |          |
|                 | ClientKeyFile             | string        | "mqtt_ssl_key_file,omitempty"        | Yes      |          |
|                 | ClientCertFile            | string        | "mqtt_ssl_cert_file,omitempty"       | Yes      |          |
|                 | ServerCaFile              | string        | "mqtt_ssl_ca_file,omitempty"         | Yes      |          |
|                 | StatsTopic                | string        | "mqtt_stats_topic,omitempty"         | Yes      |          |
| SensorConfig    | GpioPin                   | string        | "gpio_pin,omitempty"                 | Yes      |          |
|                 | GpioPollingIntervalMs     | int           | "gpio_polling_interval_ms,omitempty" | Yes      |          |


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