package internal

import "gobot.io/x/gobot"

type eventChannel chan *gobot.Event

type FakePir struct {
	Conn gobot.Connection
}

func (driver *FakePir) Name() string {
	return "FakePir"
}
func (driver *FakePir) SetName(s string) {}

func (driver *FakePir) Start() error {
	return nil
}
func (driver *FakePir) Halt() error {
	return nil
}
func (driver *FakePir) Connection() gobot.Connection {
	return driver.Conn
}

func (driver *FakePir) Pin() string {
	return "4"
}

func (driver *FakePir) Events() (eventnames map[string]string) {
	return map[string]string{}
}

func (driver *FakePir) Event(name string) string {
	return ""
}

func (driver *FakePir) AddEvent(name string)                  {}
func (driver *FakePir) DeleteEvent(name string)               {}
func (driver *FakePir) Publish(name string, data interface{}) {}

func (driver *FakePir) Unsubscribe(events eventChannel) {}
func (driver *FakePir) Subscribe() (events eventChannel) {
	return nil
}
func (driver *FakePir) On(name string, f func(s interface{})) (err error) {
	return nil
}
func (driver *FakePir) Once(name string, f func(s interface{})) (err error) {
	return nil
}
