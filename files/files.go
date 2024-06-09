// Module for finding files to lint and provide their internal representation.
package files

import (
	"io/fs"
	"os"
	"path/filepath"

	"lintex/tslatex"

	sitter "github.com/smacker/go-tree-sitter"
)

type File struct {
	Path    string
	Tree    *sitter.Node
	Source  []byte
	Ignored bool
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

func GetFiles() ([]File, error) {
	cwdFS := os.DirFS(".").(fs.ReadFileFS)
	cwd, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	paths, err := FindFiles(cwdFS, cwd)
	if err != nil {
		return nil, err
	}

	var files []File
	for _, path := range paths {
		source, err := cwdFS.ReadFile(path)
		if err != nil {
			return nil, err
		}
		tree, err := tslatex.GetTree(source)
		if err != nil {
			return nil, err
		}
		ignored, err := has_file_ignore_comment(tree, source)
		if err != nil {
			return nil, err
		}
		files = append(files, File{
			Path:    path,
			Tree:    tree,
			Source:  source,
			Ignored: *ignored,
		})
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

	result := len(matches) > 0
	return &result, nil
}
