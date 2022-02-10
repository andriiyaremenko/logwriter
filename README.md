# logw

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/andriiyaremenko/logwriter)

This package provides simple structured logging.
`logw.LogWriter(...)` returns `io.Writer` that most of the existing loggers can consume.
It will transform input according to `logw.Formatter` you choose.
You can control log level by using provided `logw.Debug|Info|Wran|Error|Fatal` variables like this:
```go
log.Println(logw.Error, "your message")
```
or like this:
```go
log.Prinln(logw.Error.WitMessage("Hello %s", "World"))
```

### To install logw:
```go
go get -u github.com/andriiyaremenko/logwriter
```

### How to use:
```go
import (
	"log"
	logw "github.com/andriiyaremenko/logwriter"
)
func main() {
	ctx := context.Background()
	log := log.New(logw.JSONLogWriter(ctx, os.Stdout), "", log.Lmsgprefix)

	log.Printf("starting work: %s", "important work")
	log.Println(logw.Warn.WithString("work", "important work"), "is done")
	// will output:
	//  {"date":"20**-**-**T**:**:**Z","level":"info","levelCode":2,"message":"starting work: important work"}
	//  {"date":"20**-**-**T**:**:**Z","level":"warn","levelCode":3,"message":"is done","work":["important work"]}
}
```
