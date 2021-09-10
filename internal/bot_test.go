package internal

/*

// NOTE: Due to the unexported type in https://github.com/hybridgroup/gobot/blob/release/eventer.go#L5 this doesn't work :-(
func TestAssembleBot(t *testing.T) {
	conf := NewConfig("test")
	mqttAdaptor := mqtt.NewAdaptor(conf.MqttConfig.Host, conf.MqttConfig.ClientId)
	fakeAdaptor := &FakeAdaptor{}
	station := &MotionDetection{
		Driver:      &FakePir{},
		Adaptor:     fakeAdaptor,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	AssembleBot(station)
}


*/
