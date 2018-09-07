package cloudlogrus

import (
	"sync"

	cl "github.com/anexia-it/go-cloudlog"
	"github.com/sirupsen/logrus"
)

// CloudlogClient interface allows you to pass your own implementation of a cloudlog client or mock clients
type CloudlogClient interface {
	PushEvent(interface{}) error
}

// Hook sends data synchronous to cloudlog
type Hook struct {
	client     CloudlogClient
	level      logrus.Level
	conFunc    func(*logrus.Entry) interface{}
	levelMutex sync.RWMutex
}

var converterFunc = func(entry *logrus.Entry) interface{} {
	if entry.Data["error"] != nil {
		errorField, ok := entry.Data["error"].(error)
		if ok {
			entry.Data["error"] = errorField.Error()
		}
	}


	d := document{
		Fields:  entry.Data,
		Message: entry.Message,
		Level:   entry.Level.String(),
	}

	return d
}

// NewHook creates a new cloudlogrus Hook to be added to logrus logger config.
func NewHook(index, ca, cert, key string) (*Hook, error) {

	client, err := cl.InitCloudLog(index, ca, cert, key)
	if err != nil {
		return nil, err
	}

	return NewCustomHook(client, converterFunc, logrus.DebugLevel), nil
}

// NewCustomHook returns a new cloudlogrus hook with a custom cloudlog client and converter
func NewCustomHook(client CloudlogClient, f func(*logrus.Entry) interface{}, l logrus.Level) *Hook {
	return &Hook{
		client:  client,
		conFunc: converterFunc,
		level:   l,
	}
}

// Must reates a new cloudlogrus Hook to be added to logrus logger config
// or panics if the hook can not be created.
func Must(index, ca, cert, key string) *Hook {
	hook, err := NewHook(index, ca, cert, key)

	if err != nil {
		panic(err)
	}
	return hook
}

type document struct {
	Message string        `cloudlog:"message"`
	Level   string        `cloudlog:"level"`
	Fields  logrus.Fields `cloudlog:"fields"`
}

// Fire sends the log entry to cloudlog
func (hook *Hook) Fire(entry *logrus.Entry) error {
	// check for custom log level filter
	if entry.Level > hook.getLevel() {
		return nil
	}

	d := hook.conFunc(entry)
	return hook.client.PushEvent(d)
}

// Levels return the available levels for this Hook. It returns logrus.AllLevels
func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// SetLevel Sets a custom log level filter for this hook
func (hook *Hook) SetLevel(l logrus.Level) {
	hook.levelMutex.Lock()
	hook.level = l
	hook.levelMutex.Unlock()
}

func (hook *Hook) getLevel() logrus.Level {
	hook.levelMutex.RLock()
	defer hook.levelMutex.RUnlock()
	return hook.level
}
