package subprocess_test

import (
	"github.com/bluele/go-subprocess"
	"testing"
)

func TestCmd(t *testing.T) {
	ret := <-subprocess.
		Cmd("ls", "-al").
		Connect()
	if ret.StatusCode != subprocess.SUCCESS {
		t.Errorf("Not expected status code %v, reason is %v.", ret.StatusCode, ret.Stderr)
	}
}

func TestCommand(t *testing.T) {
	conn := subprocess.
		Cmd("sleep", "1").
		Pipe("ls").
		WithDir("/").
		Connect()
	ret := <-conn
	if ret.StatusCode != subprocess.SUCCESS {
		t.Errorf("Not expected status code %v, reason is %v.", ret.StatusCode, ret.Stderr)
	}
}
