package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomGenerator_ErrInvalidLength(t *testing.T) {
	gen, err := NewRandom(-5)
	require.ErrorIs(t, ErrInvalidLength, err)
	require.Nil(t, gen)
}

func TestRandomGenerator_GenerateAlias(t *testing.T) {
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
			gen, err := NewRandom(tc.len)
			require.NoError(t, err)
			require.NotNil(t, gen)

			a, err := gen.NewAlias()
			require.NoError(t, err)
			require.NotNil(t, a)
			require.Len(t, a, tc.len)
		})
	}
}

func TestRandomGenerator_Uniqueness(t *testing.T) {
	gen, err := NewRandom(DefaultLength)
	require.NoError(t, err)

	a1, err := gen.NewAlias()
	require.NoError(t, err)

	a2, err := gen.NewAlias()
	require.NoError(t, err)

	require.NotEqual(t, a1, a2, "aliases must be unique")
}
