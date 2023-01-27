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
	Log(msg interface{}) error
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

// Log will print the msg to the loggers writer. It prints key/value pairs in JSON format.
func (l *logger) Log(msg interface{}) error {
	entry := map[string]interface{}{
		"msg":       msg,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
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
