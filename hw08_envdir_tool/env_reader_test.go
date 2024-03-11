package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")
		require.NoError(t, err)

		require.Equal(t, 5, len(envs))
		require.True(t, envs["UNSET"].NeedRemove)

		require.Equal(t, "bar", envs["BAR"].Value)
		require.Equal(t, `   foo
with new line`, envs["FOO"].Value)
		require.Equal(t, "", envs["EMPTY"].Value)
		require.Equal(t, "\"hello\"", envs["HELLO"].Value)
	})

	t.Run("bad dir", func(t *testing.T) {
		_, err := ReadDir("./testdata/123123")
		require.ErrorIs(t, err, ErrBadDir)
	})

	t.Run("bad file name", func(t *testing.T) {
		f, err := os.CreateTemp("./testdata/env", "sample=")
		if err != nil {
			return
		}
		defer f.Close()

		_, err = ReadDir("./testdata/env")
		os.Remove(f.Name())
		require.ErrorIs(t, err, ErrBadEnvFileName)
		fmt.Println(f.Name())
	})
}
