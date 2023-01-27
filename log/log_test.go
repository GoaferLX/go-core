package log

import (
	"bytes"
	"testing"
)

func BenchmarkLog(b *testing.B) {
	logger := New()
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Log("write something")
	}

}
