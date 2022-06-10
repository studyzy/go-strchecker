// Package strchecker finds invalid strings in code.
// comment string not check.
package strchecker

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	testSuffix = "_test.go"
)

// Parser invalid string in code parser
type Parser struct {
	path          string
	skipFile      string
	ignoreTests   bool
	excludeTypes  map[Type]bool
	invalidStrReg *regexp.Regexp
	strs          []InvalidString
}

// New creates a new instance of the parser.
// This is your entry point if you'd like to use strchecker as an API.
func New(path string, skipFile string, ignoreTests bool, excludeTypes map[Type]bool,
	invalidStrExp string) (*Parser, error) {
	//by default, ASCII only
	if len(invalidStrExp) == 0 {
		invalidStrExp = "[^\\x00-\\xff]"
	}
	invalidStrReg, err := regexp.Compile(invalidStrExp)
	if err != nil {
		return nil, err
	}
	return &Parser{
		path:          path,
		skipFile:      skipFile,
		ignoreTests:   ignoreTests,
		excludeTypes:  excludeTypes,
		invalidStrReg: invalidStrReg,
		strs:          []InvalidString{},
	}, nil
}

// ParseTree will search the given path for occurrences that could be moved into constants.
// If "..." is appended, the search will be recursive.
func (p *Parser) ParseTree() ([]InvalidString, error) {
	pathLen := len(p.path)
	// Parse recursively the given path if the recursive notation is found
	if pathLen >= 5 && p.path[pathLen-3:] == "..." {
		err := filepath.Walk(p.path[:pathLen-3],
			func(path string, f os.FileInfo, err error) error {
				if err != nil {
					log.Println(err)
					// resume walking
					return nil
				}
				if f.IsDir() {
					innerErr := p.parseDir(path)
					if innerErr != nil {
						return innerErr
					}
				}
				return nil
			})
		if err != nil {
			return nil, err
		}
	} else {
		err := p.parseDir(p.path)
		if err != nil {
			return nil, err
		}
	}

	return p.strs, nil
}

func (p *Parser) parseDir(dir string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		valid, name := true, info.Name()

		if p.ignoreTests {
			if strings.HasSuffix(name, testSuffix) {
				valid = false
			}
		}

		if len(p.skipFile) != 0 {
			match, err := regexp.MatchString(p.skipFile, dir+name)
			if err != nil {
				log.Fatal(err)
				return true
			}
			if match {
				valid = false
			}
		}

		return valid
	}, 0)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for fn, f := range pkg.Files {
			ast.Walk(&treeVisitor{
				fileSet:     fset,
				packageName: pkg.Name,
				fileName:    fn,
				p:           p,
			}, f)
		}
	}

	return nil
}

// InvalidString invalid string check result
type InvalidString struct {
	Str string
	token.Position
	packageName string
}

// Type golang code type
type Type int

// string in golang code situation
const (
	Assignment Type = iota
	Binary
	Case
	Return
	Call
	Const
)
