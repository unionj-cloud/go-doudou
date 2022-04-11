package v3

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	"strings"
	"testing"
)

func Test_isSupport(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				t: "float32",
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				t: "[]int64",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupport(tt.args.t); got != tt.want {
				t.Errorf("isSupport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_castFunc(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				t: "uint64",
			},
			want: "ToUint64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CastFunc(tt.args.t); got != tt.want {
				t.Errorf("castFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemaOf(t *testing.T) {
	Convey("SchemaOf", t, func() {
		So(SchemaOf(astutils.FieldMeta{
			Name:     "name",
			Type:     "string",
			IsExport: true,
		}), ShouldEqual, v3.String)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "age",
			Type:     "int",
			IsExport: true,
		}), ShouldEqual, v3.Int)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "id",
			Type:     "int64",
			IsExport: true,
		}), ShouldEqual, v3.Int64)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "married",
			Type:     "bool",
			IsExport: true,
		}), ShouldEqual, v3.Bool)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "score",
			Type:     "float32",
			IsExport: true,
		}), ShouldEqual, v3.Float32)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "average",
			Type:     "float64",
			IsExport: true,
		}), ShouldEqual, v3.Float64)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "avatar",
			Type:     "v3.FileModel",
			IsExport: true,
		}), ShouldEqual, v3.File)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "params",
			Type:     "...int",
			IsExport: true,
		}).Type, ShouldEqual, v3.ArrayT)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "data",
			Type:     "map[string]string",
			IsExport: true,
		}).Type, ShouldEqual, v3.ObjectT)

		So(SchemaOf(astutils.FieldMeta{
			Name:     "anony",
			Type:     "anonystruct«{\"Name\":\"\",\"Fields\":[{\"Name\":\"Name\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Name\"},{\"Name\":\"Addr\",\"Type\":\"anonystruct«{\\\"Name\\\":\\\"\\\",\\\"Fields\\\":[{\\\"Name\\\":\\\"Zip\\\",\\\"Type\\\":\\\"string\\\",\\\"Tag\\\":\\\"\\\",\\\"Comments\\\":null,\\\"IsExport\\\":true,\\\"DocName\\\":\\\"Zip\\\"},{\\\"Name\\\":\\\"Block\\\",\\\"Type\\\":\\\"string\\\",\\\"Tag\\\":\\\"\\\",\\\"Comments\\\":null,\\\"IsExport\\\":true,\\\"DocName\\\":\\\"Block\\\"},{\\\"Name\\\":\\\"Full\\\",\\\"Type\\\":\\\"string\\\",\\\"Tag\\\":\\\"\\\",\\\"Comments\\\":null,\\\"IsExport\\\":true,\\\"DocName\\\":\\\"Full\\\"}],\\\"Comments\\\":null,\\\"Methods\\\":null,\\\"IsExport\\\":false}»\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Addr\"}],\"Comments\":null,\"Methods\":null,\"IsExport\":false}»",
			IsExport: true,
		}).Type, ShouldEqual, v3.ObjectT)

		SchemaNames = []string{"User"}
		So(SchemaOf(astutils.FieldMeta{
			Name:     "user",
			Type:     "User",
			IsExport: true,
		}).Ref, ShouldEqual, "#/components/schemas/User")

		So(SchemaOf(astutils.FieldMeta{
			Name:     "user",
			Type:     "vo.User",
			IsExport: true,
		}).Ref, ShouldEqual, "#/components/schemas/User")

		Enums = map[string]astutils.EnumMeta{
			"KeyboardLayout": astutils.EnumMeta{
				Name: "KeyboardLayout",
				Values: []string{
					"UNKNOWN",
					"QWERTZ",
				},
			},
		}
		So(SchemaOf(astutils.FieldMeta{
			Name:     "layout",
			Type:     "KeyboardLayout",
			IsExport: true,
		}).Enum, ShouldResemble, []string{
			"UNKNOWN",
			"QWERTZ",
		})

		So(SchemaOf(astutils.FieldMeta{
			Name:     "any",
			Type:     "Any",
			IsExport: true,
		}), ShouldEqual, v3.Any)
	})
}

func TestRefAddDoc(t *testing.T) {
	Convey("Description should be equal to doc", t, func() {
		SchemaNames = []string{"User"}
		schema := SchemaOf(astutils.FieldMeta{
			Name:     "user",
			Type:     "User",
			IsExport: true,
		})
		doc := "This is a struct for user field"
		RefAddDoc(schema, doc)
		So(strings.TrimSpace(Schemas["User"].Description), ShouldEqual, doc)
	})
}

func TestIsBuiltin(t *testing.T) {
	Convey("Test IsBuiltIn", t, func() {
		So(IsBuiltin(astutils.FieldMeta{
			Name:     "age",
			Type:     "int",
			IsExport: true,
		}), ShouldBeTrue)

		So(IsBuiltin(astutils.FieldMeta{
			Name:     "books",
			Type:     "[]string",
			IsExport: true,
		}), ShouldBeTrue)

		So(IsBuiltin(astutils.FieldMeta{
			Name:     "data",
			Type:     "map[string]string",
			IsExport: true,
		}), ShouldBeFalse)
	})
}

func TestIsEnum(t *testing.T) {
	Convey("Test IsEnum", t, func() {
		So(IsEnum(astutils.FieldMeta{
			Name:     "age",
			Type:     "int",
			IsExport: true,
		}), ShouldBeFalse)

		Enums = map[string]astutils.EnumMeta{
			"KeyboardLayout": astutils.EnumMeta{
				Name: "KeyboardLayout",
				Values: []string{
					"UNKNOWN",
					"QWERTZ",
				},
			},
		}
		So(IsEnum(astutils.FieldMeta{
			Name:     "layout",
			Type:     "KeyboardLayout",
			IsExport: true,
		}), ShouldBeTrue)
	})
}

func TestElementType(t *testing.T) {
	Convey("Test ElementType", t, func() {
		So(ElementType("[]int"), ShouldEqual, "int")
		So(ElementType("...int"), ShouldEqual, "int")
	})
}

func TestIsOptional(t *testing.T) {
	Convey("Test IsOptional", t, func() {
		So(IsOptional("*[]int"), ShouldBeTrue)
		So(IsOptional("...int"), ShouldBeTrue)
		So(IsOptional("int"), ShouldBeFalse)
	})
}

func TestIsSlice(t *testing.T) {
	Convey("Test IsOptional", t, func() {
		So(IsSlice("*[]int"), ShouldBeTrue)
		So(IsSlice("...int"), ShouldBeTrue)
		So(IsSlice("int"), ShouldBeFalse)
	})
}
