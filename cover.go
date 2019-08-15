package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"
)

// CheckCoverage reads the coverage report specified by filename.
// For each block without coverage (count=0), the block will be analysed for a comment containing `tag`
// Any block which has no coverage and is not tagged, is recorded as an exception
func CheckCoverage(filename, tag string) error {
	failCover := false
	profiles, err := cover.ParseProfiles(filename)
	if err != nil {
		return err
	}
	for _, p := range profiles {
		fs, f, err := parseWithComments(p.FileName)
		if err != nil {
			return err // nocover
		}
		for _, block := range p.Blocks {
			if block.Count > 0 {
				continue
			}
			node := findNode(block.StartLine, block.StartCol, fs, f)
			if node == nil {
				// nocover
				return fmt.Errorf("Block not found for %s:%d.  Has the source been modified since coverage was generated?",
					p.FileName, block.StartLine)
			}
			if !isBlockTagged(f, node, tag) {
				if err = printBlock(fs, f, node); err != nil {
					return err // nocover
				}
				failCover = true
			}
		}
	}
	if failCover {
		return errors.New("FAIL: Non-tagged code exists without coverage")
	}
	return nil
}

func printBlock(fs *token.FileSet, file *ast.File, node ast.Node) error {
	config := &printer.Config{Mode: printer.SourcePos | printer.RawFormat}
	if err := config.Fprint(os.Stdout, fs, &printer.CommentedNode{Node: node, Comments: file.Comments}); err != nil {
		return err // nocover
	}
	fmt.Print("\n\n")
	return nil
}

func findNode(startLine, startCol int, fs *token.FileSet, file *ast.File) ast.Node {
	var node ast.Node
	for _, d := range file.Decls {
		ast.Inspect(d, func(n ast.Node) bool {
			if n != nil {
				pos := fs.Position(n.Pos())
				if pos.Line == startLine && pos.Column == startCol {
					node = n
					return false
				}
			}
			return true
		})
	}
	return node
}

func isBlockTagged(file *ast.File, node ast.Node, tag string) bool {
	for _, c := range file.Comments {
		if c.Pos() > node.Pos() && c.Pos() < node.End() {
			if strings.Contains(c.Text(), tag) {
				return true
			}
		}
	}
	return false
}

func parseWithComments(filename string) (*token.FileSet, *ast.File, error) {
	srcFile, err := findFile(filename)
	if err != nil {
		return nil, nil, err // nocover
	}
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, srcFile, nil, parser.ParseComments)
	return fs, f, err
}

// findFile finds the location of the named file in GOROOT, GOPATH etc.
func findFile(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err) // nocover
	}
	return filepath.Join(pkg.Dir, file), nil
}
