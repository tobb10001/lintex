package files_test

import (
	"embed"
	"io/fs"
	"lintex/files"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed assets/test_find_files/*
var find_files_fs embed.FS

func TestFindFilesFS(t *testing.T) {
	filesystem, err := fs.Sub(find_files_fs, "assets/test_find_files")
	if err != nil {
		t.Error("Couldn't get the sub-filesystem.")
	}
	files, err := files.FindFilesFS(filesystem, "")

	require.NoError(t, err)
	require.Len(t, files, 2)
	assert.Equal(t, "/directory/another_file.tex", files[0])
	assert.Equal(t, "/file.tex", files[1])
}
