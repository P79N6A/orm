package example_test

import (
	"testing"

	"github.com/op/go-logging"
)

func TestFormatFuncName(t *testing.T) {
	var log = logging.MustGetLogger("example")
	log.Error("error")
}
