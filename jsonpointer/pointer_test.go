package jsonpointer

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    Pointer
		wantErr bool
	}{
		{
			name:    "empty",
			s:       "",
			want:    Pointer{},
			wantErr: true,
		},
		{
			name:    "invalid",
			s:       "foo",
			want:    Pointer{},
			wantErr: true,
		},
		{
			name:    "root",
			s:       "/",
			want:    Pointer{},
			wantErr: false,
		},
		{
			name:    "root with token",
			s:       "/foo",
			want:    Pointer{"foo"},
			wantErr: false,
		},
		{
			name:    "root with escaped token",
			s:       "/~0",
			want:    Pointer{"~"},
			wantErr: false,
		},
		{
			name:    "multiple tokens",
			s:       "/foo/bar/baz",
			want:    Pointer{"foo", "bar", "baz"},
			wantErr: false,
		},
		{
			name:    "invalid token",
			s:       "/foo~",
			want:    Pointer{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.s, err, tt.wantErr)
				return
			}
			if err == nil && !slices.Equal(p, tt.want) {
				t.Errorf("Parse(%q) = %+v, want %+v", tt.s, p, tt.want)
			}
		})
	}
}

func TestPointer_String(t *testing.T) {
	tests := []struct {
		name string
		p    Pointer
		want string
	}{
		{
			name: "empty",
			p:    Pointer{},
			want: "/",
		},
		{
			name: "single token",
			p:    Pointer{"foo"},
			want: "/foo",
		},
		{
			name: "multiple tokens",
			p:    Pointer{"foo", "bar", "baz"},
			want: "/foo/bar/baz",
		},
		{
			name: "token with slash",
			p:    Pointer{"foo/bar"},
			want: "/foo~1bar",
		},
		{
			name: "token with tilde",
			p:    Pointer{"foo~bar"},
			want: "/foo~0bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.p.String() != tt.want {
				t.Errorf("Pointer.String() = %v, want %v", tt.p.String(), tt.want)
			}
		})
	}
}
