package logging

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	// Arrange.
	testCases := []struct {
		name     string
		logLevel Level
		in       string
		expected string
	}{
		{
			name:     "logs correctly when level >= log level",
			logLevel: InfoLevel,
			in:       "this should be logged",
			expected: "{\"level\":\"info\",\"msg\":\"this should be logged\"}\n",
		},
		{
			name:     "doesn't log when level < log level",
			logLevel: ErrorLevel,
			in:       "this should not be logged",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buf bytes.Buffer
				l   = New(&buf, Config{Level: tc.logLevel})
			)

			// Act.
			l.Info(tc.in)

			// Assert.
			assert.Equal(t, tc.expected, buf.String())
		})
	}
}

func TestInfoWith(t *testing.T) {
	// Arrange.
	var (
		buf      bytes.Buffer
		expected = `{"level":"info","msg":"test message","key1":"val1","key2":"val2"}`
	)

	l := New(&buf, Config{
		Level:         InfoLevel,
		WithTimeStamp: false, // To keep the test deterministic.
		Options:       []Option{},
	})

	// Act.
	l.InfoWith("test message",
		"key1", "val1",
		"key2", "val2",
	)

	// Assert.
	assert.JSONEq(t, expected, buf.String())
}

func TestError(t *testing.T) {
	// Arrange.
	var (
		buf      bytes.Buffer
		expected = "{\"level\":\"error\",\"msg\":\"some-error\"}\n"
	)

	l := New(&buf, Config{
		Level:         InfoLevel,
		WithTimeStamp: false, // To keep the test deterministic.
	})

	// Act.
	l.Error(fmt.Errorf("some-error"))

	// Assert.
	assert.Equal(t, expected, buf.String())
}
