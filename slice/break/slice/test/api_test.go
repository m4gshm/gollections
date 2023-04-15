package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	bslice "github.com/m4gshm/gollections/slice/break/slice"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	if homeDir, err := os.UserHomeDir(); err != nil {
		t.Error(err)
	} else if err := os.Chdir(homeDir); err != nil {
		t.Error(err)
	} else if abs, err := bslice.Convert(slice.Of("/home/user", "././inTemp"), filepath.Abs); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, slice.Of("/home/user", filepath.Join(homeDir, "inTemp")), abs)
	}
}

func TestFilterByIgnore(t *testing.T) {
	filtered, _ := bslice.Filter(slice.Of("/home/user", "./inTemp", "/usr/bin"), func(path string) error {
		return op.IfElse(filepath.IsAbs(path), nil, bslice.ErrIgnore)
	})
	assert.Equal(t, slice.Of("/home/user", "/usr/bin"), filtered)
}

func TestFilterByBreakOnFirstAppropriateAndIgnoreAnother(t *testing.T) {
	filtered, _ := bslice.Filter(slice.Of("./inTemp", "/home/user", "../siblings"), func(path string) error {
		return op.IfElse(filepath.IsAbs(path), bslice.ErrBreak, bslice.ErrIgnore)
	})
	assert.Equal(t, slice.Of("/home/user"), filtered)
}

func TestFilterByBreakOnFirstInappropriate(t *testing.T) {
	filtered, _ := bslice.Filter(slice.Of("./inTemp", "/home/user", "../siblings"), func(path string) error {
		return op.IfElse(filepath.IsAbs(path), bslice.ErrIgnoreAndBreak, nil)
	})
	assert.Equal(t, slice.Of("./inTemp"), filtered)
}
