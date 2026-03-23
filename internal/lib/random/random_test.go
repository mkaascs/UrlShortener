package random

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewRandomString(t *testing.T) {
	tests := []struct {
		name string
		len  int
	}{
		{
			name: "generate short string",
			len:  6,
		},
		{
			name: "generate long string",
			len:  128,
		},
		{
			name: "generate with negative len",
			len:  -16,
		},
		{
			name: "generate empty string",
			len:  0,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := NewRandomString(test.len)
			if test.len >= 0 {
				require.Equal(t, test.len, len(result))
			} else {
				require.Equal(t, 0, len(result))
			}
		})
	}
}
