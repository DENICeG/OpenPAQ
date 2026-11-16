// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"reflect"
	"testing"
)

func TestCHCity(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "remove space, numbers and lower case",
			input:   "City 2",
			want:    "city",
			wantErr: false,
		},
		{
			name:    "replace language specific letters",
			input:   "City íéèêëàâäùûüôöçœ",
			want:    "city ieeeeaaaeuuueooecoe",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewCh()
			got, err := g.City(tt.input)
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

func TestChPostalCode(t *testing.T) {
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
			want:    "3241",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewCh()
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

func TestChStreet(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:    "remove newline",
			input:   "the new\nroad",
			want:    []string{"the new", "road"},
			wantErr: false,
		},
		{
			name:    "remove newline with spaces",
			input:   "the new \n road",
			want:    []string{"the new", "road"},
			wantErr: false,
		},
		{
			name:    "remove newline with spaces and numbers",
			input:   "the new \n road 23",
			want:    []string{"the new", "road"},
			wantErr: false,
		},
		{
			name:    "replace str. and strasse",
			input:   "Wonderful str. strasse",
			want:    []string{"wonderful straße straße"},
			wantErr: false,
		},
		{
			name:    "replace language specific letters",
			input:   "Street íéèêëàâäùûüôöçœ",
			want:    []string{"street ieeeeaaaeuuueooecoe"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewCh()
			got, err := g.Street(tt.input)
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
