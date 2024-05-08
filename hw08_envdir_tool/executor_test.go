package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		os.Setenv("BAR", "foo")
		envs, err := ReadDir("./testdata/env")
		if err != nil {
			return
		}

		require.NoError(t, err)
		code := RunCmd([]string{"ls", "-la"}, envs)

		require.Equal(t, code, 0)
		require.Equal(t, os.Getenv("BAR"), "bar")
	})

	t.Run("error exit code", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")
		if err != nil {
			return
		}

		require.NoError(t, err)
		code := RunCmd([]string{"ls", "-23423423432"}, envs)

		require.Equal(t, code, 2)
	})
}
