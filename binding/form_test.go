package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formMultipartBinding_BindBody(t *testing.T) {
	type testObj struct {
		Param int `form:"param"`
	}
	type args struct {
		body []byte
		obj  any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test_no_error",
			args:    args{body: []byte(`param1=value1&param2=value2`), obj: make(map[string]string)},
			wantErr: false,
		},
		{
			name:    "test_parse_error",
			args:    args{body: []byte(`par;am1=value1`), obj: make(map[string]string)},
			wantErr: true,
		},
		{
			name:    "test_mapForm_error",
			args:    args{body: []byte(`param=value1`), obj: &testObj{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := formMultipartBinding{}
			err := f.BindBody(tt.args.body, tt.args.obj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_formBinding_BindBody(t *testing.T) {
	type testObj struct {
		Param int `form:"param1"`
	}
	type args struct {
		body []byte
		obj  any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test_no_error",
			args:    args{body: []byte(`param1=1&param2=value2`), obj: &testObj{}},
			wantErr: false,
		},
		{
			name:    "test_parse_error",
			args:    args{body: []byte(`par;am1=value1`), obj: make(map[string]string)},
			wantErr: true,
		},
		{
			name:    "test_mapForm_error",
			args:    args{body: []byte(`param1=value1`), obj: &testObj{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := formBinding{}
			err := f.BindBody(tt.args.body, tt.args.obj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_formPostBinding_BindBody(t *testing.T) {
	type testObj struct {
		Param int `form:"param1"`
	}
	type args struct {
		body []byte
		obj  any
	}
	tests := []struct {
		name    string
		f       formPostBinding
		args    args
		wantErr bool
	}{
		{
			name:    "test_no_error",
			args:    args{body: []byte(`param1=1&param2=value2`), obj: &testObj{}},
			wantErr: false,
		},
		{
			name:    "test_parse_error",
			args:    args{body: []byte(`par;am1=value1`), obj: make(map[string]string)},
			wantErr: true,
		},
		{
			name:    "test_mapForm_error",
			args:    args{body: []byte(`param1=value1`), obj: &testObj{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := formPostBinding{}
			err := f.BindBody(tt.args.body, tt.args.obj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
