package cloudlogrus

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"errors"
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

func TestHookWithField(t *testing.T) {
	client := MockCloudlogClient{}
	hook := NewCustomHook(&client, converterFunc, logrus.InfoLevel)
	logrus.AddHook(hook)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.WithField("key","value").Info("Message")

	require.Len(t, client.events, 1, "event count missmatch")
	event, ok := client.events[0].(document)
	require.True(t, ok)
	assert.EqualValues(t, "value", event.Fields["key"])
	assert.EqualValues(t, "Message", event.Message)
}

func TestHookWithFields(t *testing.T) {
        client := MockCloudlogClient{}
        hook := NewCustomHook(&client, converterFunc, logrus.InfoLevel)
        logrus.AddHook(hook)
        logrus.SetLevel(logrus.InfoLevel)

        logrus.WithFields(logrus.Fields{"key":"value"}).Info("Message")

		require.Len(t, client.events, 1, "event count missmatch")
        event, ok := client.events[0].(document)
        require.True(t, ok)
        assert.EqualValues(t, "value", event.Fields["key"])
        assert.EqualValues(t, "Message", event.Message)
}

func TestHookWithError(t *testing.T) {
	client := MockCloudlogClient{}
	hook := NewCustomHook(&client, converterFunc, logrus.InfoLevel)
	logrus.AddHook(hook)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.WithError(errors.New("test error")).Error("Message")

	require.Len(t, client.events, 1, "event count missmatch")
	event, ok := client.events[0].(document)
	require.True(t, ok)
	assert.EqualValues(t, "test error", event.Fields["error"])
	assert.EqualValues(t, "Message", event.Message)
	assert.EqualValues(t, "error", event.Level)
}
