package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/studyzy/go-strchecker"
)

const usageDoc = `strchecker: find invalid strings in code

Usage:

  strchecker ARGS <directory> [<directory>...]

Flags:

  -skip-file         exclude files matching the given regular expression
  -ignore-tests      exclude tests from the search (default: true)
  -output            output formatting (text or json)
  -set-exit-status   Set exit status to 2 if any issues are found
  -invalid-str       Set invalid regular expression (default: ASCII only, regular expression: [^\x00-\xff])

Examples:

  strchecker ./...
  strchecker -ignore "yacc|\.pb\." $GOPATH/src/github.com/cockroachdb/cockroach/...
  strchecker -invalid-str "[，。？！]" -output json $GOPATH/src/github.com/cockroachdb/cockroach
`

var (
	flagSkipFile      = flag.String("skip-file", "", "ignore files matching the given regular expression")
	flagIgnoreTests   = flag.Bool("ignore-tests", true, "exclude tests from the search")
	flagOutput        = flag.String("output", "text", "output formatting")
	flagInvalidStr    = flag.String("invalid-str", "", "invalid string regular expression, by default: ASCII only")
	flagInvalidType   = flag.Int("invalid-type", 0, "invalid string type: 0: ASCII only,1: ASCII and Chinese")
	flagSetExitStatus = flag.Bool("set-exit-status", false, "Set exit status to 2 if any issues are found")
)

func main() {
	flag.Usage = func() {
		usage(os.Stderr)
	}
	flag.Parse()
	log.SetPrefix("strchecker: ")

	args := flag.Args()
	if len(args) < 1 {
		usage(os.Stderr)
		os.Exit(1)
	}

	lintFailed := false
	for _, path := range args {
		anyIssues, err := run(path)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if anyIssues {
			lintFailed = true
		}
	}

	if lintFailed && *flagSetExitStatus {
		os.Exit(2)
	}
}

func run(path string) (bool, error) {
	invalidStr := *flagInvalidStr
	if *flagInvalidType == 1 { //ASCII表，中文和中文标点都允许，其他语言不允许
		invalidStr = "[^\\x00-\\xff\u4e00-\u9fa5\u3002\uff1f\uff01\uff0c\u3001\uff1b\uff1a\u201c\u201d\u2018\u2019\uff08\uff09\u300a\u300b\u3008\u3009\u3010\u3011\u300e\u300f\u300c\u300d\ufe43\ufe44\u3014\u3015\u2026\u2014\uff5e\ufe4f\uffe5]"
	}
	gco, err := strchecker.New(
		path,
		*flagSkipFile,
		*flagIgnoreTests,
		map[strchecker.Type]bool{},
		invalidStr,
	)
	if err != nil {
		return false, err
	}
	strs, err := gco.ParseTree()
	if err != nil {
		return false, err
	}
	return printOutput(strs, *flagOutput)
}

func usage(out io.Writer) {
	fmt.Fprintf(out, usageDoc)
}

func printOutput(strs []strchecker.InvalidString, output string) (bool, error) {
	switch output {
	case "json":
		jdata, err := json.Marshal(strs)
		if err != nil {
			return false, err
		}
		fmt.Println(string(jdata))
	case "text":
		for i, item := range strs {
			fmt.Printf(`%d %s:%d:%d has invalid string: "%s"`+"\n",
				i, item.Filename, item.Line, item.Column, item.Str,
			)
		}
	default:
		return false, fmt.Errorf(`Unsupported output format: %s`, output)
	}
	return len(strs) > 0, nil
}
