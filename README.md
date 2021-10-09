[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-motion-detection)](https://goreportcard.com/report/github.com/soerenschneider/gobot-motion-detection)

This project uses the [Gobot Framework](https://gobot.io/) in combination with a [PIR sensor](https://gobot.io/documentation/drivers/pir-motion-sensor/) to work as a configurable motion detection bot, being able to send simple events via MQTT and exposing machine readable metrics in the Open Metrics Format.

# MQTT Payloads

## Motion detected event
"ON"

## Motion stopped event
"OFF"

# Configuration
## Via Env Variables
| ENV                                     | Default                          | Description                                      |
|-----------------------------------------|----------------------------------|--------------------------------------------------|
| GOBOT_MOTION_DETECTION_LOCATION         | -                                | Location short name of this motion detection bot |
| GOBOT_MOTION_DETECTION_LOG_MOTIONS      | false                            | Write a log message for every motion event       |
| GOBOT_MOTION_DETECTION_GPIO_PIN         | 7                                | GPIO pin to poll                                 |
| GOBOT_MOTION_DETECTION_GPIO_POLLING_MS  | 75                               | GPIO polling frequency in milliseconds           |
| GOBOT_MOTION_DETECTION_MQTT_HOST        | gobot_motion_detection-$LOCATION | MQTT connection broker                           |
| GOBOT_MOTION_DETECTION_MQTT_CLIENT_ID   | gobot_motion_detection-$LOCATION | Client ID for the MQTT connection                |
| GOBOT_MOTION_DETECTION_MQTT_TOPIC       | gobot_motion_detection/$LOCATION | Topic to publish messages into                   |
| GOBOT_MOTION_DETECTION_METRICS_ADDR     | :9400                            | Prometheus http handler listen address           |

## Via Config File

```json
{
  "location": "location",
  "metrics_addr": ":1234",
  "gpio_pin": "7",
  "polling_interval_ms": 50,
  "log_motions": false,
  "mqtt_host": "tcp://broker:1883",
  "mqtt_client_id": "client-id",
  "mqtt_topic": "mytopic/foo"
}
```

# Metrics

This project exposes the following metrics in Open Metrics format.

| Namespace              | Subsystem | Name                               | Type    | Labels   | Help                                                              |
|------------------------|-----------|------------------------------------|---------|----------|-------------------------------------------------------------------|
| gobot_motion_detection | sensor    | motions_detected_total             | counter | location | Total amount of detected motions                                  |
| gobot_motion_detection | sensor    | motions_detected_timestamp_seconds | gauge   | location | Timestamp of the last detected motion                             |
| gobot_motion_detection | mqtt      | messages_published_total           | counter | location | The amount of published MQTT messages                             |
| gobot_motion_detection | mqtt      | message_publish_errors_total       | counter | location | Total amount of errors while trying to publish messages over MQTT |
