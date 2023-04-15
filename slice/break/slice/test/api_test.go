package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/break/predicate"
	"github.com/m4gshm/gollections/slice/break/slice"
	"github.com/stretchr/testify/assert"
)

var absPath = "/home/user"
var absPath2 = "/usr/bin"

func TestConvert(t *testing.T) {
	if homeDir, err := os.UserHomeDir(); err != nil {
		t.Error(err)
	} else if err := os.Chdir(homeDir); err != nil {
		t.Error(err)
	} else if abs, err := slice.Convert(slice.Of(absPath, "././inTemp"), filepath.Abs); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, slice.Of(absPath, filepath.Join(homeDir, "inTemp")), abs)
	}
}

func TestFilterByIgnore(t *testing.T) {
	filtered, _ := slice.Filter(slice.Of(absPath, "./inTemp", absPath2), func(path string) error {
		return op.IfElse(filepath.IsAbs(path), nil, slice.ErrIgnore)
	})
	assert.Equal(t, slice.Of(absPath, absPath2), filtered)
}

func TestFilterByBreakOnFirstAppropriateAndIgnoreAnother(t *testing.T) {
	filtered, _ := slice.Filter(slice.Of("./inTemp", absPath, "../siblings"), func(path string) error {
		return op.IfElse(filepath.IsAbs(path), slice.ErrBreak, slice.ErrIgnore)
	})
	assert.Equal(t, slice.Of(absPath), filtered)
}

func TestFilterByBreakOnFirstInappropriate(t *testing.T) {
	filtered, _ := slice.Filter(slice.Of("./inTemp", absPath, "../siblings"), func(path string) error {
		return op.IfElse(filepath.IsAbs(path), slice.ErrIgnoreAndBreak, nil)
	})
	assert.Equal(t, slice.Of("./inTemp"), filtered)
}

func TestGroup(t *testing.T) {
	var (
		paths      = slice.Of("./inTemp", absPath, "../siblings", absPath2)
		grouped, _ = slice.Group(paths, predicate.Of(filepath.IsAbs))
	)
	assert.Equal(t, slice.Of(absPath, absPath2), grouped[true])
	assert.Equal(t, 2, len(grouped[false]))
}
