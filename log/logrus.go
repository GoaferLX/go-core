package log

import (
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

type Logger interface {
	Fatal(msg ...interface{})
	Fatalf(format string, msg ...interface{})
	Fatalln(msg ...interface{})

	Panic(msg ...interface{})
	Panicf(format string, msg ...interface{})
	Panicln(msg ...interface{})

	Print(msg ...interface{})
	Printf(format string, msg ...interface{})
	Println(msg ...interface{})
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

// ErrorLogger returns a Logger that can be used for HTTP servers.
func ErrorLogger() (*log.Logger, func()) {
	w := logrus.StandardLogger().Writer()
	errLog := log.New(w, "", 0)
	teardown := func() {
		if err := w.Close(); err != nil {
			WithError(err).Error("ErrorLogger pipewriter failed to close")
			return
		}
		logrus.Info("ErrorLogger pipewriter closed")
	}
	return errLog, teardown
}

// Error logs a message at level Error.
func Error(msg ...interface{}) {
	logrus.Error(msg)
}

// Fatal logs a message at level Fatal.
func Fatal(msg ...interface{}) {
	logrus.Fatal(msg)
}

func Fatalf(format string, msg ...interface{}) {
	logrus.Fatalf(format, msg)
}

// Info logs a message at level Info.
func Info(msg ...interface{}) {
	logrus.Info(msg)
}

// Print is an alias for Info.
// Print exists for compatibility reasons.  Any code relying on stdlib
// log package can import this package and it will still work without any changes.
func Print(msg ...interface{}) {
	logrus.Print(msg)
}

// Printf is an alias for Printf.
// Printf exists for compatibility reasons.  Any code relying on stdlib
// log package can import this package and it will still work without any changes.
func Printf(format string, args ...interface{}) {
	logrus.Printf(format, args)
}

// WithError prepares a log entry with the value of the error added under the 'error' key.
func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

// WithField prepares a log entry with the value added under the field key.
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(fields))
}

// HTTPRequest is a helper to log information from an incoming HTTP Request.
func HTTPRequest(r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"protocol": r.Proto,
		"method":   r.Method,
		"URI":      r.RequestURI,
	}).Info("HTTP request")
}
