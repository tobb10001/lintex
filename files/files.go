// Module for finding files to lint and provide their internal representation.
package files

import (
	"io/fs"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
)

type File struct {
	absolute_path string
	tree *sitter.Node
}

func FindFilesFS(filesystem fs.FS, prefix string) ([]string, error) {
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
	for i := range files {
		files[i] = filepath.Clean(prefix + "/" + files[i])
	}
	return files, nil
}

func FindFiles() ([]string, error) {
	fs := os.DirFS(".")	
	cwd, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	return FindFilesFS(fs, cwd)
}
