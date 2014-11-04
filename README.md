# Go-SubProcess

This is a convenience wrapper around the `os/exec` module.

# Getting Started

## Install

```
go get github.com/bluele/go-subprocess
```

## Usage

```go
// Get channel for subprocess.
result := <-subprocess.
  Cmd("ls", "-al").
  Pipe("wc", "-l").
  Connect()

// $ ls -al | wc -l
fmt.Println(result.Stdout)
```

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>
