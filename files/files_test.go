package files_test

import (
	"embed"
	"io/fs"
	"testing"

	"lintex/files"

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
	files, err := files.FindFiles(filesystem, "")

	require.NoError(t, err)
	require.Len(t, files, 2)
	assert.Equal(t, "directory/another_file.tex", files[0])
	assert.Equal(t, "file.tex", files[1])
}

func TestFileIgnore(t *testing.T) {
	for _, testcase := range []struct {
		Name    string
		Source  []byte
		Ignored bool
	}{
		{
			Name:    "standard comment",
			Source:  []byte("% lintex: ignore_file"),
			Ignored: true,
		},
		{
			Name:    "stuff after space",
			Source:  []byte("% lintex: ignore_file some stuff after space"),
			Ignored: true,
		},
		{
			Name:    "comment without lintex_ignore",
			Source:  []byte("% there is something else in this comment"),
			Ignored: false,
		},
		{
			Name:    "no comment at all",
			Source:  []byte("There is no comment."),
			Ignored: false,
		},
		{
			Name:    "no space after lintex_ignore",
			Source:  []byte("% lintex: ignore_filebut somebody forgot the space after."),
			Ignored: false,
		},
	} {
		t.Run(testcase.Name, func(t *testing.T) {
			f, err := files.NewFile("dummy_path", testcase.Source)

			require.NoError(t, err)
			assert.Equal(t, testcase.Ignored, f.Ignored())
		})
	}
}
