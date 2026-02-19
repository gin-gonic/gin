package binding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomDecimalUnmarshalParam(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "leading dot",
			input:   ".1",
			want:    "0.1",
			wantErr: false,
		},
		{
			name:    "invalid decimal",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "leading dot with multiple digits",
			input:   ".123",
			want:    "0.123",
			wantErr: false,
		},
		{
			name:    "normal decimal",
			input:   "1.23",
			want:    "1.23",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cd CustomDecimal
			err := cd.UnmarshalParam(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, cd.String())
		})
	}
}
