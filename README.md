# go-strchecker
golang code string checker [【中文】](README_CN.md)

strchecker can scan golang source code, and find all strings that matched "invalid-str" regular expression. 
By default, strchecker can find all not ASCII code string.
# Get Started
```bash
go install github.com/studyzy/go-strchecker/cmd/strchecker@latest
strchecker ./...
```
# Usage
```bash
strchecker ARGS <directory> [<directory>...]

Flags:

  -skip-file         exclude files matching the given regular expression
  -ignore-tests      exclude tests from the search (default: true)
  -output            output formatting (text or json)
  -set-exit-status   Set exit status to 2 if any issues are found
  -invalid-str       Set invalid regular expression (default: ASCII only, regular expression: [^\x00-\xff])
```

# Examples
```bash
strchecker ./...
strchecker -skip-file "_mock.go" $GOPATH/src/github.com/studyzy/iocgo
strchecker -invalid-str "[，。？！]" -output json $GOPATH/src/github.com/studyzy/iocgo
strchecker -invalid-type=1 ./testdata/...
```
# Result
```bash
0 testdata/main.go:10:30 has invalid string: "not found！"
1 testdata/main.go:12:17 has invalid string: "no，data！"
2 testdata/main.go:15:14 has invalid string: "Hello，World！"
3 testdata/main.go:16:12 has invalid string: "Current time："
4 testdata/main.go:19:15 has invalid string: "한국어"
5 testdata/main.go:20:15 has invalid string: "にほんご"
6 testdata/main.go:22:14 has invalid string: ":) 😁😁😁"
7 testdata/call.go:9:60 has invalid string: "！"
8 testdata/call.go:10:11 has invalid string: "a！b"
9 testdata/call.go:11:5 has invalid string: "aa！"
10 testdata/call.go:12:40 has invalid string: "bb！"
```