// Module for finding files to lint and provide their internal representation.
package files

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"lintex/tslatex"

	sitter "github.com/smacker/go-tree-sitter"
)

type File struct {
	absolute_path string
	tree          *sitter.Node
	source        []byte
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

func GetFiles() ([]File, error) {
	paths, err := FindFiles()
	if err != nil {
		return nil, err
	}

	var files []File
	for _, path := range paths {
		source, err := Read(path)
		if err != nil {
			return nil, err
		}
		tree, err := tslatex.GetTree(source)
		if err != nil {
			return nil, err
		}
		files = append(files, File{
			absolute_path: path,
			tree:          tree,
			source:        source,
		})
	}
	return files, nil
}
