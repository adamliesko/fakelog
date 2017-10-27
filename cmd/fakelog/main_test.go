package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/adamliesko/fakelog"
)

func TestRun(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "fakelog")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	//tmp file, 10 requests per second, 1 second
	err = run(filepath.Join(f.Name()), "common", 10, 1)
	if err != nil {
		t.Error(err)
	}

	gotLines := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		gotLines = true
		if l := scanner.Text(); l == "" {
			t.Error("produced empty log line", l)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if !gotLines {
		t.Fatalf("nothing was written to the file")
	}

}

func TestRunNonExistentFile(t *testing.T) {
	//dir as file, 10 requests per second, 1 second
	err := run(os.TempDir(), "", 10, 1)
	if err == nil {
		t.Error("shouldn't be able to write to '' file")
	}
}

func TestRunWithStdout(t *testing.T) {
	//Stdout as file, 10 requests per second, 1 second
	err := run("", "combined", 5, 1)
	if err != nil {
		t.Errorf("should be able to write to stdout, but got an error: %v", err)
	}
}

func TestRunUnknownLogFormat(t *testing.T) {
	err := run("", "log-format-of-my-custom-imaginary-toolkit", 10, 1)
	if err == nil {
		t.Error("should have got an error with unknown log format")
	}
}

func TestPickingLineGenerator(t *testing.T) {
	tcs := []struct {
		in  string
		out fakelog.LineGenerator
		err error
	}{
		{
			in:  "",
			out: fakelog.ApacheCommonLine,
		},
		{
			in:  "common",
			out: fakelog.ApacheCommonLine,
		},
		{
			in:  "combined",
			out: fakelog.ApacheCombinedLine,
		},
		{
			in:  "nginx-custom-logger",
			err: errors.New("unknown log format: nginx-custom-logger"),
		},
	}
	for _, tc := range tcs {
		fn, err := pickLineGenerator(tc.in)
		if tc.out != nil && getFunctionName(fn) != getFunctionName(tc.out) {
			t.Errorf("wrong line generator, got %v want %v", getFunctionName(fn), getFunctionName(tc.out))
		}
		if (tc.err != nil && (tc.err.Error() != err.Error())) || (tc.err == nil && err != nil) {
			t.Errorf("returned error mismatch, got '%v' want '%v'", err, tc.err)
		}
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
