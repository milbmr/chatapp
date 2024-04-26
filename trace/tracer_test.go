package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
  var buf bytes.Buffer
  tracer := New(&buf)
  if tracer == nil {
    t.Error("tracer shouldn't contain nil")
  } else {
    tracer.Trace("hello world.")
    if buf.String() != "hello world.\n" {
      t.Errorf("trace shouldn't write %s", buf.String())
    }
  }
}

func TestOff(t *testing.T) {
  var silentTracer Tracer = Off()
  silentTracer.Trace("something")
}
