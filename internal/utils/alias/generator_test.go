package alias

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerator_ErrInvalidLength(t *testing.T) {
	gen, err := New(-5)
	require.ErrorIs(t, ErrInvalidLength, err)
	require.Nil(t, gen)
}

func TestGenerator_GenerateAlias(t *testing.T) {
	cases := []struct {
		name string
		len  int
	}{
		{"len 5", 5},
		{"len 10", 10},
		{"len 100", 100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gen, err := New(tc.len)
			require.NoError(t, err)
			require.NotNil(t, gen)

			a, err := gen.NewAlias()
			require.NoError(t, err)
			require.NotNil(t, a)
			require.Len(t, a, tc.len)
		})
	}
}
