package cloudlogrus

import (
	"testing"

	"github.com/sirupsen/logrus"
)

type MockCloudlogClient struct {
	events []interface{}
}

func (client *MockCloudlogClient) PushEvent(e interface{}) error {
	client.events = append(client.events, e)
	return nil
}

var testLevels = []struct {
	level    logrus.Level
	expected int
}{
	{logrus.DebugLevel, 4},
	{logrus.InfoLevel, 3},
	{logrus.WarnLevel, 2},
	{logrus.ErrorLevel, 1},
}

func TestHook(t *testing.T) {
	client := MockCloudlogClient{}
	hook := NewCustomHook(&client, converterFunc, logrus.InfoLevel)
	logrus.AddHook(hook)
	logrus.SetLevel(logrus.DebugLevel)

	for _, e := range testLevels {
		client.events = nil
		hook.SetLevel(e.level)

		logrus.Debug("")
		logrus.Info("")
		logrus.Warn("")
		logrus.Error("")

		l := len(client.events)
		if l != e.expected {
			t.Errorf("len is %v, expected %v", l, e.expected)
		}
	}
}
