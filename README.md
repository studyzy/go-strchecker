# go-strchecker
golang code string checker
# Get Started
> go install github.com/studyzy/go-strchecker/cmd/strchecker@latest
> strchecker ./...
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
```