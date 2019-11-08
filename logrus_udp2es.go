package logrus_udp2es

import (
	"os"
	"net"
	"fmt"
	"time"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type connectInterface interface {
	Write([]byte) (int, error)
}

type Hook struct {
	// Connection Details
	Host string
	Port int

	// es index
	ESIndex string

	levels []logrus.Level

	conn connectInterface
}

// NewPapertrailHook creates a UDP hook to be added to an instance of logger.
func NewUdp2EsHook(hook *Hook) (*Hook, error) {
	var err error
	hook.conn, err = net.Dial("udp", fmt.Sprintf("%s:%d", hook.Host, hook.Port))

	return hook, err
}

// Fire is called when a log envent is fired.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	logStr, _ := entry.String()

	var logDetail map[string]interface{}
	if err := json.Unmarshal([]byte(logStr), &logDetail); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to unmarshal log string: %s", logStr)
		return err
	}

	logDetail["time"] = time.Now().Unix() * 1000
	logDetail["level"] = entry.Level.String()
	logDetail["index"] = hook.ESIndex
	if message, ok := logDetail["msg"]; ok {
		logDetail["message_analyzed"] = message
	}

	var payload []byte
	payload, err := json.Marshal(logDetail)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal json log detail: %+v", logDetail)
		return err
	}

	bytesWritten, err := hook.conn.Write(payload)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to send log line to udp2es. Wrote %d bytes before error: %v", bytesWritten, err)
		return err
	}

	return nil
}

// SetLevels specify nessesary levels for this hook.
func (hook *Hook) SetLevels(lvs []logrus.Level) {
	hook.levels = lvs
}

// Levels returns the available logging levels.
func (hook *Hook) Levels() []logrus.Level {
	if hook.levels == nil {
		return []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		}
	}

	return hook.levels
}
