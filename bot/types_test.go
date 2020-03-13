package bot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommand_Name(t *testing.T) {
	t.Run("", func(t *testing.T) {
		cmd := NewCommand("test", nil)

		assert.Equal(t, "test", cmd.Name())
		assert.Equal(t, map[string]interface{}{}, cmd.Args())
	})

	assert.PanicsWithValue(t, `command name "test data" is wrong, can't contain spaces or be an empty string`, func() {
		NewCommand("test data", nil)
	})

	assert.PanicsWithValue(t, `command name "" is wrong, can't contain spaces or be an empty string`, func() {
		NewCommand("", nil)
	})

	assert.PanicsWithValue(t, `command name "   " is wrong, can't contain spaces or be an empty string`, func() {
		NewCommand("   ", nil)
	})
}

func TestCommand_Match(t *testing.T) {
	assert.True(t, NewCommand("test", nil).Match("test"))
	assert.True(t, NewCommand("test", nil).Match("test val1 val2"))
	assert.False(t, NewCommand("test", nil).Match("test1"))
	assert.False(t, NewCommand("test", nil).Match(" test"))
}
