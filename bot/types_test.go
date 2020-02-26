package bot

import (
	"github.com/JILeXanDR/skypebot/skypeapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommand_Name(t *testing.T) {
	tt := []struct {
		cmd          Command
		expectedName string
		expectedArgs map[string]interface{}
	}{
		{Command("test"), "test", map[string]interface{}{}},
		{Command("test data"), "test", map[string]interface{}{}},
		{Command(""), "", map[string]interface{}{}},
		{Command("   "), "", map[string]interface{}{}},
	}

	for _, tc := range tt {
		tc := tc
		assert.Equal(t, tc.expectedName, tc.cmd.Name(), "name")
		assert.Equal(t, tc.expectedArgs, tc.cmd.Args(&Activity{activity: &skypeapi.Activity{Text: "test"}}), "name")
	}
}

func TestCommand_Match(t *testing.T) {
	assert.True(t, Command("test").Match("test"))
	assert.True(t, Command("test").Match("test val1 val2"))
	assert.False(t, Command("test").Match("test1"))
}
