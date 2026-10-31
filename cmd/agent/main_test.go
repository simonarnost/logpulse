package main

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/simonarnost/logpulse/gen/logpulsepb"
)

// fakeSender records sent entries; returns err from Send if set.
type fakeSender struct {
	sent []*logpulsepb.LogEntry
	err  error
}

func (f *fakeSender) Send(e *logpulsepb.LogEntry) error {
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, e)
	return nil
}

func messages(entries []*logpulsepb.LogEntry) []string {
	var out []string
	for _, e := range entries {
		out = append(out, e.Message)
	}
	return out
}

func TestStreamLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty input", "", nil},
		{"single line no newline", "one", []string{"one"}},
		{"multiple lines", "one\ntwo\nthree\n", []string{"one", "two", "three"}},
		{"skips empty lines", "one\n\n\ntwo\n", []string{"one", "two"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakeSender{}
			if err := streamLines(strings.NewReader(tt.input), f); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := messages(f.sent); !slices.Equal(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStreamLinesSendError(t *testing.T) {
	sendErr := errors.New("send failed")
	f := &fakeSender{err: sendErr}
	if err := streamLines(strings.NewReader("one\n"), f); !errors.Is(err, sendErr) {
		t.Errorf("got %v, want %v", err, sendErr)
	}
}

func TestStreamLinesReadError(t *testing.T) {
	readErr := errors.New("read failed")
	f := &fakeSender{}
	if err := streamLines(iotest.ErrReader(readErr), f); !errors.Is(err, readErr) {
		t.Errorf("got %v, want %v", err, readErr)
	}
}
