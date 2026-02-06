package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate_URL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want error
	}{
		{"success http", "http://example.com", nil},
		{"success https", "https://example.com/path?Q=1", nil},
		{"success ssh", "ssh://example.com", nil},

		{"failed", "", ErrInvalidURL},
		{"failed", "   ", ErrInvalidURL},
		{"failed", "example.com", ErrInvalidURL},
		{"failed", "ht!tp://example.com", ErrInvalidURL},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := URL(tc.in)
			if tc.want == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, tc.want, err)
			}
		})
	}
}
