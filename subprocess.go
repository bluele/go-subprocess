package subprocess

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"time"
)

const (
	SUCCESS = 0
	FAILED  = 1
	TIMEOUT = 2
)

type SubProcess struct {
	commands []*Command

	env     []string
	dir     string
	stdin   io.Reader
	timeout time.Duration
}

type Command struct {
	name string
	args []string
}

func Cmd(name string, arg ...string) *SubProcess {
	sp := SubProcess{}
	sp.commands = append(sp.commands, &Command{
		name: name,
		args: arg,
	})
	return &sp
}

func (sp *SubProcess) Pipe(name string, arg ...string) *SubProcess {
	sp.commands = append(sp.commands, &Command{
		name: name,
		args: arg,
	})
	return sp
}

func (sp *SubProcess) SetEnv(key, value string) *SubProcess {
	sp.env = append(sp.env, key+"="+value)
	return sp
}

func (sp *SubProcess) WithDir(dir string) *SubProcess {
	sp.dir = dir
	return sp
}

func (sp *SubProcess) WithTimeout(to time.Duration) *SubProcess {
	sp.timeout = to
	return sp
}

func (sp *SubProcess) SetStdin(ir io.Reader) *SubProcess {
	sp.stdin = ir
	return sp
}

type Result struct {
	Stdout     io.Reader
	Stderr     io.Reader
	StatusCode int
}

func (sp *SubProcess) Connect() chan *Result {
	ch := make(chan *Result)
	go func() {
		var ret *Result
		if sp.timeout == time.Duration(0) {
			ret = <-sp.connect()
		} else {
			select {
			case ret = <-sp.connect():
			case <-time.After(sp.timeout):
				ret = createTimeoutResult()
			}
		}
		ch <- ret
	}()
	return ch
}

func (sp *SubProcess) connect() chan *Result {
	ch := make(chan *Result)
	go func() {
		ch <- sp.execCommand()
	}()
	return ch
}

func (sp *SubProcess) execCommand() *Result {
	var out bytes.Buffer
	var ec *exec.Cmd
	var stdout io.ReadCloser
	var isHead bool = true
	var err error
	var ecs []*exec.Cmd
	var length = len(sp.commands)

	for _, cmd := range sp.commands {
		ec = exec.Command(cmd.name, cmd.args...)
		if isHead {
			ec.Stdin = sp.stdin
			isHead = false
		} else {
			ec.Stdin = stdout
		}
		ec.Env = sp.env
		ec.Dir = sp.dir

		stdout, err = ec.StdoutPipe()
		if err != nil {
			return &Result{
				Stderr:     strings.NewReader(err.Error()),
				StatusCode: FAILED,
			}
		}
		ecs = append(ecs, ec)
	}

	for i := length; i > 0; i-- {
		ec = ecs[i-1]
		if i == length {
			ec.Stdout = &out
		}
		ec.Start()
	}
	for i := 0; i < length; i++ {
		ecs[i].Wait()
	}

	return &Result{
		Stdout:     &out,
		StatusCode: SUCCESS,
	}
}

func createTimeoutResult() *Result {
	return &Result{
		Stderr:     strings.NewReader("Connection timeout error"),
		StatusCode: TIMEOUT,
	}
}
