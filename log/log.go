/*
package log aims to provide a simple logging library that can be used on its own or in conjunction with existing logging libraries.
The main purpose of this package is to remove the fluff provided by other packages, with a smaller and more succint API,
whilst allowing the caller to leverage those packages if required.
*/
package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Logger defines a simple interface to be passed through applications.  It does exactly what it says on the tin, logs.
// This reduces the API of a logger and its ability to affect an application once it is in use,
// making its usage simpler to understand/use, whilst avoiding unintended side effects.
type Logger interface {
	Log(msg string, fields ...interface{}) error
}

// logger is a very simple implementation of the  Logger interface.
// It provides structured json logging by default.
type logger struct {
	writer io.Writer  // destinatation for log output.
	mu     sync.Mutex //  mutex prevents concurrent writes to the output.
}

func New() *logger {
	return &logger{
		writer: os.Stdout,
		mu:     sync.Mutex{},
	}
}

// DefaultLogger exposes a pre-configfured implementation of the Logger interface that can used when the caller does not require more fine grained control over
// their logger instance.  Similar to the stdlib pattern.
var DefaultLogger *logger = New()

// Log will print the msg to the loggers writer. It prints key/value pairs in JSON format.
// Fields are key/value pairs that will be logged to provide additional ocntext.  If there is an odd number of pairs, they will be silently dropped.
func (l *logger) Log(msg string, fields ...any) error {
	entry := map[string]string{
		"msg":       msg,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	if len(fields)%2 == 0 {
		for i := 0; i < len(fields); {
			fieldName, ok := fields[i].(string)
			if ok {
				entry[fieldName] = fmt.Sprint(fields[i+1])
			}
			i = i + 2
		}
	}

	bytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("log: marshalling json: %w", err)
	}

	l.mu.Lock()
	// write the output, each entry is on a new line by default.
	_, err = l.writer.Write(append(bytes, byte('\n')))
	l.mu.Unlock()

	return err
}

// SetOutput exists so that the log destination can be safely changed whilst being protected
// by the mutex.  Originally this was not going to be included, as it seems YAGNI, however,
// it was necessary for testing.
func (l *logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	l.writer = w
	l.mu.Unlock()
}
