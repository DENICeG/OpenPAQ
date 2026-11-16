// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"reflect"
	"testing"
)

func TestGeneric_City(t *testing.T) {

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "remove space, numbers and lower case",
			args:    args{s: "Theodor-Stern-Kai 2"},
			want:    "theodor-stern-kai",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewGeneric()
			got, err := g.City(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("City() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("City() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneric_PostalCode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "remove space and lower case",
			args:    args{s: "3241 ML"},
			want:    "3241ml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewGeneric()
			got, err := g.PostalCode(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("City() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("City() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneric_Street(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "remove newline",
			args:    args{s: "the new\nroad"},
			want:    []string{"the new", "road"},
			wantErr: false,
		},
		{
			name:    "remove newline with spaces",
			args:    args{s: "the new \n road"},
			want:    []string{"the new", "road"},
			wantErr: false,
		},
		{
			name:    "remove newline with spaces and numbers",
			args:    args{s: "the new \n road 23"},
			want:    []string{"the new", "road"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewGeneric()
			got, err := g.Street(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Street() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Street() got = %v, want %v", got, tt.want)
			}
		})
	}
}
