package copier

import (
	"fmt"
	"testing"
)

type Family struct {
	Father string
	Mather string
	Pets   map[string]string
	Toys   []string
}

type TestStruct struct {
	Name   string
	Age    int
	Family Family
}

type FamilyShadow struct {
	Father string
	Mather string
	Pets   map[string]string
}

type TestStructShadow struct {
	Name   string
	Family FamilyShadow
}

func TestDeepCopy(t *testing.T) {
	pets := make(map[string]string)
	pets["a"] = "dog"
	pets["b"] = "cat"

	family := Family{
		Father: "Jack",
		Mather: "Lily",
		Pets:   pets,
		Toys: []string{
			"car",
			"lego",
		},
	}
	src := TestStruct{
		Name:   "Rose",
		Age:    18,
		Family: family,
	}

	var target TestStructShadow

	type args struct {
		src    interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				src:    src,
				target: &target,
			},
			wantErr: false,
		},
		{
			args: args{
				src:    nil,
				target: nil,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				src:    make(chan string),
				target: &target,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeepCopy(tt.args.src, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("DeepCopy() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Printf("%#v\n", tt.args.target)
		})
	}
}

func TestDeepCopy_ShouldHasError(t *testing.T) {
	pets := make(map[string]string)
	pets["a"] = "dog"
	pets["b"] = "cat"

	family := Family{
		Father: "Jack",
		Mather: "Lily",
		Pets:   pets,
		Toys: []string{
			"car",
			"lego",
		},
	}
	src := TestStruct{
		Name:   "Rose",
		Age:    18,
		Family: family,
	}

	var target TestStructShadow

	type args struct {
		src    interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestDeepCopy",
			args: args{
				src:    src,
				target: target,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeepCopy(tt.args.src, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("DeepCopy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
