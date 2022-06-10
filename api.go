package strchecker

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

// Issue api call this module, return invalid string check list as issue
type Issue struct {
	Pos         token.Position
	Str         string
	PackageName string
	Index       int
}

// Config api call config
type Config struct {
	IgnoreTests   bool
	invalidStrReg string
	SkipFile      string
	ExcludeTypes  map[Type]bool
}

// Run api run this module, check every input go source files
func Run(files []*ast.File, fset *token.FileSet, cfg *Config) ([]Issue, error) {
	p, err := New(
		"",
		cfg.SkipFile,
		cfg.IgnoreTests,
		cfg.ExcludeTypes,
		cfg.invalidStrReg,
	)
	if err != nil {
		return nil, err
	}
	var skipFileReg *regexp.Regexp
	if len(p.skipFile) > 0 {
		skipFileReg, err = regexp.Compile(p.skipFile)
		if err != nil {
			return nil, err
		}
	}
	var issues []Issue
	for _, f := range files {
		if p.ignoreTests {
			if filename := fset.Position(f.Pos()).Filename; strings.HasSuffix(filename, testSuffix) {
				continue
			}
		}
		if skipFileReg != nil {
			if filename := fset.Position(f.Pos()).Filename; skipFileReg.MatchString(filename) {
				continue
			}
		}
		ast.Walk(&treeVisitor{
			fileSet:     fset,
			packageName: "",
			fileName:    "",
			p:           p,
		}, f)
	}

	for i, item := range p.strs {
		issue := Issue{
			Pos:         item.Position,
			Str:         item.Str,
			PackageName: item.packageName,
			Index:       i,
		}
		issues = append(issues, issue)
	}
	return issues, nil
}
