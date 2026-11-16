// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"reflect"
	"testing"
)

func TestUK_City(t *testing.T) {

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
			name:    "remove invalid characters",
			args:    args{s: "city name +/(){}[]<>!§$%&=?*#€¿_\",:;\n\n"},
			want:    "city name",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewGB()
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

func TestUK_PostalCode(t *testing.T) {
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
			name:    "valid postal code 1",
			args:    args{s: "A1 6CD"},
			want:    "a1",
			wantErr: false,
		}, {
			name:    "valid postal code 2",
			args:    args{s: "AA1 6CD"},
			want:    "aa1",
			wantErr: false,
		}, {
			name:    "valid postal code 3",
			args:    args{s: "A12 6CD"},
			want:    "a12",
			wantErr: false,
		}, {
			name:    "valid postal code 4",
			args:    args{s: "AA12 6CD"},
			want:    "aa12",
			wantErr: false,
		}, {
			name:    "valid postal code 5",
			args:    args{s: "A1B 6CD"},
			want:    "a1b",
			wantErr: false,
		}, {
			name:    "valid postal code 6",
			args:    args{s: "AA1B 6CD"},
			want:    "aa1b",
			wantErr: false,
		},
		{
			name:    "valid postal code with space at start and end",
			args:    args{s: " AA1B 6CD "},
			want:    "aa1b",
			wantErr: false,
		},
		{
			name:    "valid postal code with space at start and end and missing middle space",
			args:    args{s: " AA1B6CD "},
			want:    "aa1b",
			wantErr: false,
		}, {
			name:    "invalid postal code 1",
			args:    args{s: " AA1B6C "},
			want:    "",
			wantErr: true,
		}, {
			name:    "invalid postal code 2",
			args:    args{s: " AA1B6C "},
			want:    "",
			wantErr: true,
		}, {
			name:    "invalid postal code 3",
			args:    args{s: " AA1C "},
			want:    "",
			wantErr: true,
		}, {
			name:    "invalid postal code 4",
			args:    args{s: " 1AA1C "},
			want:    "",
			wantErr: true,
		}, {
			name:    "stripe special characters",
			args:    args{s: " AA12 6CD +/(){}[]<>!§'$%&=?*#€¿_\":;\n\n- "},
			want:    "aa12",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uk, _ := NewGB()
			got, err := uk.PostalCode(tt.args.s)
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

func TestUK_Street(t *testing.T) {
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
			name:    "must be lowercase",
			args:    args{s: "The New One"},
			want:    []string{"the new one"},
			wantErr: false,
		},
		{
			name:    "remove invalid characters",
			args:    args{s: "The New One +./(){}[]<>!§'$%&=?*#€¿_\":;0123456789-\n"},
			want:    []string{"the new one"},
			wantErr: false,
		},
		{
			name:    "remove postal code",
			args:    args{s: "The New One AB1 1AB AB11AB"},
			want:    []string{"the new one"},
			wantErr: false,
		},
		{
			name:    "addresses contain road or street",
			args:    args{s: "The New Street, ANd the new ROad"},
			want:    []string{"new street", "the new street", "the new road", "and the new road"},
			wantErr: false,
		},
		{
			name:    "addresses split from number till street",
			args:    args{s: "What ever 123 The New Street AB1 1AB Some City"},
			want:    []string{"the new street", "what ever the new street some city"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, _ := NewGB()
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
