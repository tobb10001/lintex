// Module for finding files to lint and provide their internal representation.
package files

import (
	"io/fs"
	"os"
	"path/filepath"

	"lintex/tslatex"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

type File struct {
	path    string
	tree    *sitter.Node
	source  []byte
	ignored bool
}

func NewFile(path string, source []byte) (*File, error) {
	tree, err := tslatex.GetTree(source)
	if err != nil {
		return nil, err
	}
	ignored, err := has_file_ignore_comment(tree, source)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("source", string(source)).Bool("ignored", *ignored).Send()
	return &File{
		path:    path,
		tree:    tree,
		source:  source,
		ignored: *ignored,
	}, nil
}

func (f *File) Ignored() bool  { return f.ignored }
func (f *File) Path() string   { return f.path }
func (f *File) Source() []byte { return f.source }

func (f *File) GetMatches(pattern []byte) (*sitter.Query, []*sitter.QueryMatch, error) {
	return tslatex.GetMatches(f.tree, pattern, f.source)
}

func FindFiles(filesystem fs.FS, prefix string) ([]string, error) {
	var files []string
	error := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".tex" {
			files = append(files, path)
		}
		return nil
	})
	if error != nil {
		return nil, error
	}
	return files, nil
}

func GetFiles() ([]*File, error) {
	cwdFS := os.DirFS(".").(fs.ReadFileFS)
	cwd, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	paths, err := FindFiles(cwdFS, cwd)
	if err != nil {
		return nil, err
	}

	var files []*File
	for _, path := range paths {
		source, err := cwdFS.ReadFile(path)
		if err != nil {
			return nil, err
		}
		currentFile, err := NewFile(path, source)
		if err != nil {
			return nil, err
		}
		files = append(files, currentFile)
	}
	return files, nil
}

func has_file_ignore_comment(tree *sitter.Node, source []byte) (*bool, error) {
	pattern := []byte(`
		(
			(line_comment) @comment
			(#match? @comment "\\% lintex: ignore_file($| )")
		)
	`)

	_, matches, err := tslatex.GetMatches(tree, pattern, source)
	if err != nil {
		return nil, err
	}

	log.Debug().Int("matches", len(matches)).Send()

	result := len(matches) > 0
	return &result, nil
}
