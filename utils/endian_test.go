package utils

import "testing"

func TestBigToLittleEndian(t *testing.T) {
	type test struct {
		name        string
		input       string
		want        string
		expectedErr bool
	}
	tests := []test{
		{
			name:        "valid hex string",
			input:       "fefb",
			want:        "fbfe",
			expectedErr: false,
		},
		{
			name:        "invalid hex string",
			input:       "fezz",
			want:        "",
			expectedErr: true,
		},
		{
			name:        "valid long hex string",
			input:       "fefbac1b7d2101d0",
			want:        "d001217d1bacfbfe",
			expectedErr: false,
		},
	}
	for _, tc := range tests {
		out, err := BigToLittleEndian(tc.input)
		if err == nil && tc.expectedErr {
			t.Fatalf("%s failed, expected error but didn't get one", tc.input)
		}
		if err != nil && !tc.expectedErr {
			t.Fatalf("%s failed, got error [%v] and didn't expect one", tc.input, err)
		}
		if out != tc.want {
			t.Fatalf("%s failed, wanted [%v], got [%v]", tc.input, tc.want, out)
		}
	}
}
