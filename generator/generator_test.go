package generator

import (
	"errors"
	"regexp"
	"testing"
	"time"
)

const (
	ApacheCommonLineRegex   = `^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+)\s?(\S+)?\s?(\S+)?" (\d{3}|-) (\d+|-)$`
	ApacheCombinedLineRegex = `^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+)\s?(\S+)?\s?(\S+)?" (\d{3}|-) (\d+|-)\s?"?([^"]*)"?\s?"?([^"]*)?"?$`
)

type testOut struct{ count int }

func (to *testOut) Write(p []byte) (n int, err error) {
	to.count++
	return 0, nil
}
func (to *testOut) Close() error { return nil }

type testOutError struct {
	testOut
}

func (toe *testOutError) Write(p []byte) (n int, err error) {
	toe.count++
	return 0, errors.New("full disk")
}

func TestFakeLoggerWithTooBigRate_GenerateLogs(t *testing.T) {
	to := &testOut{}
	fl := NewFakeLogger(ApacheCommonLine, to, 20000)

	if fl.rate > maxRate {
		t.Fatalf("rate %d it bigger then the allowed maxRate %d", fl.rate, maxRate)
	}
}

func TestErroneousWriteCloser(t *testing.T) {
	to := &testOutError{}
	fl := NewFakeLogger(ApacheCommonLine, to, 20000)

	err := fl.GenerateLogs()
	if err == nil {
		t.Fatalf("didn't get an error")
	}
	if got := to.count; got != 1 {
		t.Fatalf("got too many attempted writes, generator should have stopped after 1, got %d", got)
	}
}

func TestFakeLogger_Stop(t *testing.T) {
	to := &testOut{}
	fl := NewFakeLogger(ApacheCommonLine, to, 200)

	// run the logger async and wait it produces some logs
	go func() {
		err := fl.GenerateLogs()
		if err != nil {
			t.Fatal(err)
		}
	}()

	start := time.Now().Add(1 * time.Second)
	for {
		if to.count > 50 {
			break
		}
		if time.Now().After(start) {
			fl.Stop()
			t.Fatalf("slow to produce log lines, %d lines produced", to.count)
		}
		time.Sleep(10 * time.Millisecond)
	}
	fl.Stop()
	// record current count of written lines
	written := to.count
	// sleep a little to allow FakeLogger to potentially keep logging
	time.Sleep(1 * time.Second)

	if written != to.count {
		t.Errorf("Stop didn't stop the fake logger from producing logs, got %v want %v", to.count, written)
	}

}

func TestApacheCommonLogLine(t *testing.T) {
	reg, err := regexp.Compile(ApacheCommonLineRegex)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		l := ApacheCommonLine()
		if !reg.MatchString(l) {
			t.Errorf("produced line '%s' doesn't follow apache common log format", l)
		}
	}
}

func TestApacheCombinedLogLine(t *testing.T) {
	reg, err := regexp.Compile(ApacheCombinedLineRegex)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		l := ApacheCombinedLine()
		if !reg.MatchString(l) {
			t.Errorf("produced line '%s' doesn't follow apache combined log format", l)
		}
	}
}

func BenchmarkGenerator_ApacheCommonLogLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ApacheCommonLine()
	}
}

func BenchmarkGenerator_ApacheCombinedLogLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ApacheCombinedLine()
	}
}
