// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"reflect"
	"testing"
)

func TestFrCity(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "remove space, numbers and lower case",
			input:   " City   +/(){}[]<>!§$%&=?*#€¿_\",:;023456789\n\n",
			want:    "city",
			wantErr: false,
		},
		{
			name:    "remove space, numbers and lower case",
			input:   " City íéèêëàâäùûüôöçœ",
			want:    "city ieeeeaaauuuoocoe",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := newFR()
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

func Test_PostalCode_fr(t *testing.T) {
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
			args:    args{s: "3241 Ml"},
			want:    "3241",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := newFR()
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

func TestFrStreet(t *testing.T) {

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
			name:    "remove newline with spaces and numbers",
			input:   "The New \n road 23",
			want:    []string{"the new", "road"},
			wantErr: false,
		},
		{
			name:    "Split at -",
			input:   "The New - road 23",
			want:    []string{"the new", "road"},
			wantErr: false,
		},

		{
			name:    "Remove - at start and",
			input:   "The New -   -road- ",
			want:    []string{"the new", "road"},
			wantErr: false,
		},
		{
			name:    "replace language specific letters",
			input:   "street íéèêëàâäùûüôöçœ",
			want:    []string{"street ieeeeaaauuuoocoe"},
			wantErr: false,
		},
		{
			name:    "replace bd.",
			input:   "bd. asdf",
			want:    []string{"boulevard asdf"},
			wantErr: false,
		},
		{
			name:    "replace bd",
			input:   "bd asdf",
			want:    []string{"boulevard asdf"},
			wantErr: false,
		},
		{
			name:    "replace ave",
			input:   "ave asdf",
			want:    []string{"avenue asdf"},
			wantErr: false,
		},
		{
			name:    "replace zac",
			input:   "zac qwery ave asdf ",
			want:    []string{"zac qwery avenue asdf", "avenue asdf"},
			wantErr: false,
		},
		{
			name:    "replace zac",
			input:   "ave asdf zac qwery",
			want:    []string{"avenue asdf", "zac qwery"},
			wantErr: false,
		},

		{
			name:    "append avenue",
			input:   "vorname ave asdf",
			want:    []string{"vorname avenue asdf", "avenue asdf"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := newFR()
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
