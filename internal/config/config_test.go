package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_fromEnvBool1(t *testing.T) {
	key := "asdjnasdogsagsadgjsdgsdgasdgjsdg"
	resultingKey := fmt.Sprintf("%s_%s", strings.ToUpper(BotName), strings.ToUpper(key))
	os.Setenv(resultingKey, "true")
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				name: "asdjfasdighasgasdgasdg",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "default",
			args: args{
				name: "asdjfasdighasgasdgasdg",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "test",
			args: args{
				name: key,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromEnvBool(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromEnvBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fromEnvBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{
			name:    "no tld",
			host:    "tcp://hostname:1883",
			wantErr: false,
		},
		{
			name:    "tld",
			host:    "tcp://hostname.my.tld:1883",
			wantErr: false,
		},
		{
			name:    "ip",
			host:    "tcp://192.168.0.1:1883",
			wantErr: false,
		},
		{
			name:    "no protocol",
			host:    "192.168.0.1:1883",
			wantErr: true,
		},
		{
			name:    "no port",
			host:    "tcp://host",
			wantErr: true,
		},
		{
			name:    "only host",
			host:    "host",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := matchHost(tt.host); (err != nil) != tt.wantErr {
				t.Errorf("matchHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fromEnvInt(t *testing.T) {
	key := "akjsdfjasgasdkg"
	resultingKey := fmt.Sprintf("%s_%s", strings.ToUpper(BotName), strings.ToUpper(key))
	os.Setenv(resultingKey, "3141")

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    key,
			want:    3141,
			wantErr: false,
		},
		{
			name:    "bla",
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromEnvInt(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromEnvInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fromEnvInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Placement             string
		MetricConfig          string
		FirmAtaPort           string
		GpioPin               string
		GpioPollingIntervalMs int
		IntervalSecs          int
		LogValues             bool
		MqttConfig            MqttConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all okay",
			fields: fields{
				Placement:             "Placement",
				MetricConfig:          ":9100",
				FirmAtaPort:           "/dev/ttyUSB0",
				GpioPin:               "5",
				GpioPollingIntervalMs: 75,
				IntervalSecs:          30,
				LogValues:             false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: false,
		},
		{
			name: "missing loc",
			fields: fields{
				MetricConfig:          ":9100",
				FirmAtaPort:           "/dev/ttyUSB0",
				GpioPin:               "5",
				GpioPollingIntervalMs: 75,
				IntervalSecs:          30,
				LogValues:             false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
		{
			name: "missing gpiopin",
			fields: fields{
				Placement:             "loc",
				MetricConfig:          ":9100",
				GpioPollingIntervalMs: 75,
				IntervalSecs:          30,
				LogValues:             false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
		{
			name: "missing host",
			fields: fields{
				Placement:             "loc",
				MetricConfig:          ":9100",
				GpioPin:               "5",
				GpioPollingIntervalMs: 75,
				IntervalSecs:          30,
				LogValues:             false,
				MqttConfig: MqttConfig{
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Placement:   tt.fields.Placement,
				MetricsAddr: tt.fields.MetricConfig,
				SensorConfig: SensorConfig{
					GpioPin:               tt.fields.GpioPin,
					GpioPollingIntervalMs: tt.fields.GpioPollingIntervalMs,
				},
				IntervalSecs: tt.fields.IntervalSecs,
				LogSensor:    tt.fields.LogValues,
				MqttConfig:   tt.fields.MqttConfig,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadJsonConfig(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     *Config
		wantErr  bool
	}{
		{
			name:     "non-existent-file",
			filePath: "ihopethispathdoesntexist/somefile.json",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "example-config",
			filePath: "../../contrib/example-config.json",
			want: &Config{
				Placement:   "loc",
				MetricsAddr: defaultMetricConfig,
				SensorConfig: SensorConfig{
					GpioPin:               defaultGpioPin,
					GpioPollingIntervalMs: defaultGpioPollingIntervalMs,
				},
				IntervalSecs: defaultIntervalSeconds,
				LogSensor:    defaultLogValues,
				MqttConfig: MqttConfig{
					Host:       "tcp://host:1883",
					Topic:      "sensors/pir",
					StatsTopic: "sensors/pir/stats",
				},
				StatIntervals: []int{1, 2, 3},
				MessageOn:     defaultMessageOn,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadJsonConfig(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJsonConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadJsonConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchTopic(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		wantErr bool
	}{
		{
			topic:   "topicname",
			wantErr: false,
		},
		{
			topic:   "more/complicated",
			wantErr: false,
		},
		{
			topic:   "more/complicated/topic",
			wantErr: false,
		},
		{
			topic:   "/leading",
			wantErr: true,
		},
		{
			topic:   "trailing/",
			wantErr: true,
		},
		{
			topic:   "replace/%s",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := matchTopic(tt.topic); (err != nil) != tt.wantErr {
				t.Errorf("matchTopic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_TemplateTopic(t *testing.T) {
	type fields struct {
		Placement  string
		MqttConfig MqttConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   *Config
	}{
		{
			name: "template",
			fields: fields{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix/%s",
				},
			},
			want: &Config{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix/loc",
				},
			},
		},
		{
			name: "no templating",
			fields: fields{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix",
				},
			},
			want: &Config{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Config{
				Placement:  tt.fields.Placement,
				MqttConfig: tt.fields.MqttConfig,
			}
			conf.FormatTopic()
			if !reflect.DeepEqual(conf, tt.want) {
				t.Fail()
			}
		})
	}
}
