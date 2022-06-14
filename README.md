# miilog

A [zap](https://github.com/uber-go/zap) wrapper for a slacker :zany_face:

## Example

```go
package main

import (
	"github.com/miiton/miilog"
)

func init() {
	miilog.SetLoggerProductionMust()
	// or miilog.SetLoggerDevelopmentMust()
	// or miilog.SetLoggerProductionWithLokiMust("https://loki.example.com/", "MYTENANTID", "{source=\"mygo_source\", job=\"mygo_job\", host=\"awesomehost\"}")
}

func main() {
	defer miilog.Sync()
	miilog.Infow("hoge",
		"job", "fuga",
		"moreinfo", "moremore",
	)
	miilog.Warn("ahhhhhhhhh?")
	miilog.Errorf("%s", "omgomgomgomgomg")
}
```

### Output

`SetLoggerProductionMust()`

```
{"level":"info","ts":"2022-06-15T06:54:43.072224552+09:00","caller":"tmp-zap/main.go:15","msg":"hoge","job":"fuga","moreinfo":"moremore"}
{"level":"warn","ts":"2022-06-15T06:54:43.072297729+09:00","caller":"tmp-zap/main.go:19","msg":"ahhhhhhhhh?"}
{"level":"error","ts":"2022-06-15T06:54:43.072306135+09:00","caller":"tmp-zap/main.go:20","msg":"omgomgomgomgomg","stacktrace":"main.main\n\t/home/USERNAME/dev/src/github.com/miiton/tmp-zap/main.go:20\nruntime.main\n\t/home/linuxbrew/.linuxbrew/Cellar/go/1.18/libexec/src/runtime/proc.go:250"}
```

`SetLoggerDevelopmentMust()`

```
2022-06-15T06:55:24.467+0900    INFO    tmp-zap/main.go:15      hoge    {"job": "fuga", "moreinfo": "moremore"}
2022-06-15T06:55:24.467+0900    WARN    tmp-zap/main.go:19      ahhhhhhhhh?
2022-06-15T06:55:24.467+0900    ERROR   tmp-zap/main.go:20      omgomgomgomgomg
main.main
        /home/USERNAME/dev/src/github.com/miiton/tmp-zap/main.go:20
runtime.main
        /home/linuxbrew/.linuxbrew/Cellar/go/1.18/libexec/src/runtime/proc.go:250
```
