package tslatex_test

import (
	"context"
	"testing"

	"lintex/tslatex"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	n, err := sitter.ParseCtx(context.Background(), []byte(`\documentclass{report}`), tslatex.GetLanguage())
	assert.NoError(err)
	assert.Equal(
		"(source_file (class_include path: (curly_group_path path: (path))))",
		n.String(),
	)
}
