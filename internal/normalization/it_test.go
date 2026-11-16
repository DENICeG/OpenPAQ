// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"reflect"
	"testing"
)

func TestIt_City(t *testing.T) {

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
		{
			name:    "remove invalid chars",
			args:    args{s: "Theodor-Stern-Kai 2 +/(){}[]<>!§$%&=?*#€¿_\",:;013456789\n\n"},
			want:    "theodor-stern-kai",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewIT()
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

func TestIt_PostalCode(t *testing.T) {
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
			name:    "remove space, characters and lower case",
			args:    args{s: "some messy text 12345 ML"},
			want:    "12345",
			wantErr: false,
		},
		{
			name:    "invalid postal code",
			args:    args{s: "some messy text 1234 ML"},
			want:    "1234",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewIT()
			got, err := g.PostalCode(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostalCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PostalCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIt_Street(t *testing.T) {
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
		{
			name:    "replace short name",
			args:    args{s: "the new \n c.so 23"},
			want:    []string{"the new", "corso"},
			wantErr: false,
		},
		{
			name:    "replace short name 2",
			args:    args{s: "the new \n l.go 23"},
			want:    []string{"the new", "largo"},
			wantErr: false,
		},
		{
			name:    "replace short name 3",
			args:    args{s: "the new trav. 23"},
			want:    []string{"the new traversa"},
			wantErr: false,
		},
		{
			name:    "replace short name 4",
			args:    args{s: "p.za del torro  23"},
			want:    []string{"piazza del torro"},
			wantErr: false,
		},
		{
			name:    "street include piazza",
			args:    args{s: "bullshit piazza del torro 123"},
			want:    []string{"bullshit piazza del torro", "piazza del torro"},
			wantErr: false,
		},
		{
			name:    "street include via",
			args:    args{s: "bullshit via del torro 123"},
			want:    []string{"bullshit via del torro", "via del torro"},
			wantErr: false,
		},
		{
			name: "street include g.",
			args: args{s: "via g. torro 123"},
			want: []string{
				"via g. torro",
				"via giovanni torro",
				"via giuseppe torro",
				"via giacomo torro",
				"via gabriele torro",
				"via giorgio torro",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewIT()
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
