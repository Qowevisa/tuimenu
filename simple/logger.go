package simple

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// BufferedLogger structure buffer log messages
type BufferedLogger struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

// Logf buffers the log message with a formatted string
// To print it instantly call Flush
func (bl *BufferedLogger) Logf(format string, args ...interface{}) {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	// Create log entry with a timestamp
	timestamp := time.Now().Format(time.RFC3339)
	entry := fmt.Sprintf("[%s] %s", timestamp, fmt.Sprintf(format, args...))

	if !(strings.Contains(entry, "\n") && entry[len(entry)-1] == '\n') {
		entry += "\n"
	}
	bl.buffer.WriteString(entry)
}

// Flush outputs all buffered log messages
func (bl *BufferedLogger) Flush() {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	if bl.buffer.Len() > 0 {
		fmt.Print(bl.buffer.String())
		bl.buffer.Reset()
	}
}

func (m *Menu) RedirectLogOutputToBufferedLogger() {
	log.SetOutput(&m.Log.buffer)
}
