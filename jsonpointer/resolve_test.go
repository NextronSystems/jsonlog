package jsonpointer

import "testing"

type BaseStruct struct {
	Field1     string      `json:"field1"`
	SliceField []string    `json:"slice_field"`
	SubStruct  *BaseStruct `json:"sub_struct"`
}

func TestResolve(t *testing.T) {
	var base = BaseStruct{
		Field1:     "field1",
		SliceField: []string{"slice_field"},
		SubStruct: &BaseStruct{
			Field1: "sub_field1",
		},
	}
	for _, tt := range []struct {
		desc    string
		pointer string
		want    any
	}{
		{
			desc:    "simple field",
			pointer: "/field1",
			want:    &base.Field1,
		},
		{
			desc:    "slice field",
			pointer: "/slice_field/0",
			want:    &base.SliceField[0],
		},
		{
			desc:    "sub struct field",
			pointer: "/sub_struct/field1",
			want:    &base.SubStruct.Field1,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			pointer, err := Parse(tt.pointer)
			if err != nil {
				t.Fatalf("Parse(%q) = %v, want nil", tt.pointer, err)
			}
			got, err := Resolve(&base, pointer)
			if err != nil {
				t.Errorf("Resolve(%q) = %v, want nil", tt.pointer, err)
			}
			if got != tt.want {
				t.Errorf("Resolve(%q) = %q, want %q", tt.pointer, got, tt.want)
			}
		})
	}
}
