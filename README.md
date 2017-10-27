
# fakelog - Fake log generation
[![Build Status](https://secure.travis-ci.org/adamliesko/fakelog.svg)](http://travis-ci.org/adamliesko/fakelog)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamliesko/fakelog)](https://goreportcard.com/report/github.com/adamliesko/fakelog)
[![GoDoc](https://godoc.org/github.com/adamliesko/fakelog?status.svg)](https://godoc.org/github.com/adamliesko/fakelog) 
[![Coverage Status](https://img.shields.io/coveralls/adamliesko/fakelog.svg)](https://coveralls.io/r/adamliesko/fakelog?branch=master)


`fakelog` is a package and cmd-line tool for fake and random log
 generation. It supports writing the logs to Stdout, files or buffers,
 accepting io.WriteCloser interface. Also, support for user defined
 functions generating log lines is implemented.

## Installation

```
go get github.com/adamliesko/fakelog
```

## Usage

Generate approximately 1 rps for 5 seconds to STDOUT.

```
$ ./fakelog  --rate=2 -duration=5
  5.122.23.165 sarah_cooper - [15/10/2017:19:50:31 +0200] "PUT /articles HTTP/1.1" 401 357
  154.242.0.43 abrv - [15/10/2017:19:50:31 +0200] "GET /article/5504 HTTP/1.1" 200 4052
  99.139.225.19 leet_coder - [15/10/2017:19:50:32 +0200] "GET /article/7982 HTTP/1.1" 301 919
  235.105.29.223 grace_hooper - [15/10/2017:19:50:32 +0200] "GET /login HTTP/1.1" 503 16572
  29.245.71.97 grace_hooper - [15/10/2017:19:50:33 +0200] "GET /users HTTP/1.1" 500 12429
  89.246.77.246 sarah_cooper - [15/10/2017:19:50:33 +0200] "GET /articles HTTP/1.1" 403 9786
  111.86.180.40 leet_coder - [15/10/2017:19:50:34 +0200] "PATCH /popular HTTP/1.1" 301 23116
  30.192.115.4 leet_coder - [15/10/2017:19:50:34 +0200] "DELETE /article/323 HTTP/1.1" 401 6303
  229.146.184.20 leet_coder - [15/10/2017:19:50:35 +0200] "GET /popular HTTP/1.1" 200 25138
  12.48.221.25 john_doe - [15/10/2017:19:50:35 +0200] "PUT /articles HTTP/1.1" 200 7602
```

Generate infinite number of log lines into a file, with the default rate of requests using apache combined format.
```
$ ./fakelog -format=combined -path=apache.log
```

 Providing custom function generating log lines.
```
import (
    "time"
    "github.com/adamliesko/fakelog/generator"
 )

logFn := func() string { return time.Now().String() + " hey" }
fl := fakelog.Logger(logFn, os.Stdout, 200)
fl.GenerateLogs()

Output:
2017-10-17 21:05:56.026488 +0200 CEST hey
2017-10-17 21:05:56.032415 +0200 CEST hey
2017-10-17 21:05:56.039889 +0200 CEST hey
...
```

## Help
```
$ fakelog -h
Usage of ./fakelog:
  -duration int
    	duration [s] to produce logs, 0 is inf
  -format string
    	log format (default "common")
  -path string
    	path to the log file (default os.Stdout)
  -rate int
    	approximate rate of requests per second, max of [rate, 10^6] (default 50)
```