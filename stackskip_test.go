package stackdriver_test

import (
	"bytes"
	"encoding/json"
	"testing"

	stackdriver "github.com/g5search/logrus-stackdriver-formatter"
	"github.com/g5search/logrus-stackdriver-formatter/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestStackSkip(t *testing.T) {
	var out bytes.Buffer

	logger := logrus.New()
	logger.Out = &out
	logger.Formatter = stackdriver.NewFormatter(
		stackdriver.WithService("test"),
		stackdriver.WithVersion("0.1"),
		stackdriver.WithStackSkip("github.com/g5search/logrus-stackdriver-formatter/test"),
		stackdriver.WithNoTimestamp(),
	)

	mylog := test.LogWrapper{
		Logger: logger,
	}

	mylog.Error("my log entry")

	want := map[string]interface{}{
		"severity": "ERROR",
		"message":  "my log entry",
		"serviceContext": map[string]interface{}{
			"service": "test",
			"version": "0.1",
		},
		"context": map[string]interface{}{
			"reportLocation": map[string]interface{}{
				"file":     "github.com/g5search/logrus-stackdriver-formatter/stackskip_test.go",
				"line":     30.0,
				"function": "TestStackSkip",
			},
		},
		"sourceLocation": map[string]interface{}{
			"file":     "github.com/g5search/logrus-stackdriver-formatter/stackskip_test.go",
			"line":     30.0,
			"function": "TestStackSkip",
		},
	}

	wantedBytes, err := json.Marshal(want)
	if err != nil {
		t.Error(err)
	}
	assert.JSONEq(t, string(wantedBytes), out.String())
}
