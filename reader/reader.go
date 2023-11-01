// Package to provide functionality needed to read files.
package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"lintex/tslatex"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/rs/zerolog/log"
)

type File struct {
	Path string
	Tree *sitter.Node
	Source []byte
}

// Higher level method to read a file.
func Read(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return bs, nil
}

func ReadDocument(mainpath string) ([]File, []string, error) {
	filelist := NewFileList(mainpath)

	result := []File{}
	notFound := []string{}

	for filelist.Next() {
		file := filelist.Get()
		log.Trace().Msgf("Reading file %s", file)

		source, err := Read(file)
		if errors.Is(err, os.ErrNotExist){
			notFound = append(notFound, file)
			continue
		} else if err != nil {
			return nil, nil, err
		}

		tree, includes, err := ReadFile(source)
		if err != nil {
			return nil, nil, err
		}

		result = append(result, File{Path: file, Tree: tree, Source: source})
		filelist.Add(includes...)
	}
	return result, notFound, nil
}

// Read a file, parse it, return the tree and paths to included files.
func ReadFile(source []byte) (*sitter.Node, []string, error) {
	tree, err := tslatex.GetTree(source)
	if err != nil {
		return nil, nil, err
	}

	included, err := GetIncludedFiles(tree, source)
	if err != nil {
		return nil, nil, err
	}

	return tree, included, nil
}

// Find the file paths for files, that are included in the given tree.
func GetIncludedFiles(tree *sitter.Node, source []byte) ([]string, error) {
	inputPattern := []byte(`
		(latex_include
			(curly_group_path
				(path) @path	
			)
		)
	`)

	query, matches, err := tslatex.GetMatches(tree, inputPattern, source)
	if err != nil {
		return nil, err
	}
	results := []string{}
	for _, match := range matches {
		for _, capture := range match.Captures {
			if query.CaptureNameForId(capture.Index) == "path" {
				file := capture.Node.Content(source)
				if !strings.HasSuffix(file, ".tex") {
					file += ".tex"
				}
				results = append(results, file)
			}
		}
	}
	return results, nil
}
