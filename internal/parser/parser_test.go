package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseString_BasicKeyValues(t *testing.T) {
	input := `
DB_HOST=localhost
DB_PORT=5432
APP_NAME=myapp
`
	env, err := ParseString(input)
	require.NoError(t, err)
	assert.Equal(t, "localhost", env["DB_HOST"])
	assert.Equal(t, "5432", env["DB_PORT"])
	assert.Equal(t, "myapp", env["APP_NAME"])
}

func TestParseString_IgnoresComments(t *testing.T) {
	input := `
# This is a comment
KEY=value
# Another comment
`
	env, err := ParseString(input)
	require.NoError(t, err)
	assert.Len(t, env, 1)
	assert.Equal(t, "value", env["KEY"])
}

func TestParseString_IgnoresEmptyLines(t *testing.T) {
	input := "\n\nKEY=value\n\n"
	env, err := ParseString(input)
	require.NoError(t, err)
	assert.Len(t, env, 1)
}

func TestParseString_StripQuotes(t *testing.T) {
	input := `
SINGLE='hello world'
DOUBLE="hello world"
NONE=hello
`
	env, err := ParseString(input)
	require.NoError(t, err)
	assert.Equal(t, "hello world", env["SINGLE"])
	assert.Equal(t, "hello world", env["DOUBLE"])
	assert.Equal(t, "hello", env["NONE"])
}

func TestParseString_InvalidSyntax(t *testing.T) {
	input := "INVALID_LINE_NO_EQUALS"
	_, err := ParseString(input)
	assert.Error(t, err)
}

func TestParseString_ValueWithEquals(t *testing.T) {
	input := "URL=http://example.com?foo=bar"
	env, err := ParseString(input)
	require.NoError(t, err)
	assert.Equal(t, "http://example.com?foo=bar", env["URL"])
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	assert.Error(t, err)
}

func TestParseFile_ValidFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "*.env")
	require.NoError(t, err)
	defer os.Remove(tmp.Name())

	_, err = tmp.WriteString("FOO=bar\nBAZ=qux\n")
	require.NoError(t, err)
	tmp.Close()

	env, err := ParseFile(tmp.Name())
	require.NoError(t, err)
	assert.Equal(t, "bar", env["FOO"])
	assert.Equal(t, "qux", env["BAZ"])
}
