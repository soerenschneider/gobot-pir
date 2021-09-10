package internal

type FakeAdaptor struct {
}

func (driver *FakeAdaptor) Name() string {
	return "FakeAdaptor"
}

func (driver *FakeAdaptor) SetName(n string) {
}

func (driver *FakeAdaptor) Connect() error {
	return nil
}
func (driver *FakeAdaptor) Finalize() error {
	return nil
}
