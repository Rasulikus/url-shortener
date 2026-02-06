package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounterGenerator_ErrInvalidLength(t *testing.T) {
	gen, err := NewCounter(0, 123, -10)
	require.ErrorIs(t, ErrInvalidLength, err)
	require.Nil(t, gen)
}

func TestCounterGenerator_GenerateAlias_Length(t *testing.T) {
	cases := []struct {
		name   string
		length int
	}{
		{"len 5", 5},
		{"len 10", 10},
		{"len 100", 100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gen, err := NewCounter(0, 123456, tc.length)
			require.NoError(t, err)
			require.NotNil(t, gen)

			a, err := gen.NewAlias()
			require.NoError(t, err)
			require.Len(t, a, tc.length)
		})
	}
}

func TestCounterGenerator_Uniqueness(t *testing.T) {
	gen, err := NewCounter(0, 42, DefaultLength)
	require.NoError(t, err)

	a1, err := gen.NewAlias()
	require.NoError(t, err)

	a2, err := gen.NewAlias()
	require.NoError(t, err)

	require.NotEqual(t, a1, a2, "aliases must be unique for different ids")
}
