package stackdriver_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	stackdriver "github.com/g5search/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	for _, tt := range formatterTests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer

			logger := logrus.New()
			logger.Out = &got
			logger.Formatter = stackdriver.NewFormatter(
				stackdriver.WithService("test"),
				stackdriver.WithVersion("0.1"),
				stackdriver.WithNoTimestamp(),
			)

			tt.run(logger)
			out, err := json.Marshal(tt.out)
			if err != nil {
				t.Error(err)
			}
			assert.JSONEq(t, string(out), got.String())
		})
	}
}

var formatterTests = []struct {
	run  func(*logrus.Logger)
	out  map[string]interface{}
	name string
}{
	{
		name: "With Field",
		run: func(logger *logrus.Logger) {
			logger.WithField("foo", "bar").Info("my log entry")
		},
		out: map[string]interface{}{
			"severity": "INFO",
			"message":  "my log entry",
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
				"reportLocation": map[string]interface{}{
					"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
					"line":     45,
					"function": "glob..func1",
				},
			},
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"sourceLocation": map[string]interface{}{
				"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
				"line":     45,
				"function": "glob..func1",
			},
		},
	},
	{
		name: "WithField and WithError",
		run: func(logger *logrus.Logger) {
			logger.
				WithField("foo", "bar").
				WithError(errors.New("test error")).
				Info("my log entry")
		},
		out: map[string]interface{}{
			"severity": "INFO",
			"message":  "my log entry: test error",
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo":   "bar",
					"error": "test error",
				},
				"reportLocation": map[string]interface{}{
					"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
					"line":     77,
					"function": "glob..func2",
				},
			},
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"sourceLocation": map[string]interface{}{
				"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
				"line":     77,
				"function": "glob..func2",
			},
		},
	},
	{
		name: "WithField and Error",
		run: func(logger *logrus.Logger) {
			logger.WithField("foo", "bar").Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
				"reportLocation": map[string]interface{}{
					"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
					"line":     107,
					"function": "glob..func3",
				},
			},
			"sourceLocation": map[string]interface{}{
				"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
				"line":     107,
				"function": "glob..func3",
			},
		},
	},
	{
		name: "WithField, WithError and Error",
		run: func(logger *logrus.Logger) {
			logger.
				WithField("foo", "bar").
				WithError(errors.New("test error")).
				Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry: test error",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo":   "bar",
					"error": "test error",
				},
				"reportLocation": map[string]interface{}{
					"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
					"line":     139,
					"function": "glob..func4",
				},
			},
			"sourceLocation": map[string]interface{}{
				"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
				"line":     139,
				"function": "glob..func4",
			},
		},
	},
	{
		name: "WithField, HTTPRequest and Error",
		run: func(logger *logrus.Logger) {
			logger.
				WithFields(logrus.Fields{
					"foo": "bar",
					"httpRequest": map[string]interface{}{
						"requestMethod": "GET",
					},
				}).
				Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"serviceContext": map[string]interface{}{
				"service": "test",
				"version": "0.1",
			},
			"context": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
					"httpRequest": map[string]interface{}{
						"requestMethod": "GET",
					},
				},
				"reportLocation": map[string]interface{}{
					"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
					"line":     176,
					"function": "glob..func5",
				},
			},
			"sourceLocation": map[string]interface{}{
				"file":     "github.com/g5search/logrus-stackdriver-formatter/formatter_test.go",
				"line":     176,
				"function": "glob..func5",
			},
		},
	},
}
