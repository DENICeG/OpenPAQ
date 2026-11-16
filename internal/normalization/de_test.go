// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"reflect"
	"regexp"
	"testing"
)

func TestCheckAndConvertPlz(t *testing.T) {
	de, _ := NewDE()
	type args struct {
		plz string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Telekom", args: struct{ plz string }{plz: "D-12345"}, want: "12345", wantErr: false},
		{name: "Telekom-zu-kurz", args: struct{ plz string }{plz: "D-124"}, want: "", wantErr: true},
		{name: "D-Zero-plz", args: struct{ plz string }{plz: "D-00000"}, want: "", wantErr: true},
		{name: "Zero-Plz", args: struct{ plz string }{plz: "00000"}, want: "", wantErr: true},
		{name: "Usual Case", args: struct{ plz string }{plz: "12345"}, want: "12345", wantErr: false},
		{name: "Prefix", args: struct{ plz string }{plz: "DE12345"}, want: "12345", wantErr: false},
		{name: "String-in-PLZ", args: struct{ plz string }{plz: "DE12Blub345"}, want: "12345", wantErr: false},
		{name: "Forgot-leading-0", args: struct{ plz string }{plz: "1067"}, want: "01067", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := de.PostalCode(tt.args.plz)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAndConvertPlz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAndConvertPlz() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanCity(t *testing.T) {
	de, _ := NewDE()
	type args struct {
		city string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "ö-to-oe", args: args{city: "Köln"}, want: "koeln", wantErr: false},
		{name: "ß-to-sz", args: args{city: "Gießen"}, want: "gieszen", wantErr: false},
		{name: "ü-to-ue", args: args{city: "Rüsselsheim"}, want: "ruesselsheim", wantErr: false},
		{name: "replace-", args: args{city: "Eine-Stadt"}, want: "eine stadt", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := de.City(tt.args.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("cleanCity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cleanCity() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_detectQuadratAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{name: "mannheim-deluxe", args: struct {
			address string
		}{address: "G 3 4"}, want: "g3", want1: true},
		{name: "mannheim-deluxe-different-writing", args: struct {
			address string
		}{address: "G3 4"}, want: "g3", want1: true},
		{name: "mannheim-no-housenumber", args: struct {
			address string
		}{address: "G3"}, want: "g3", want1: true},
		{name: "Fake Quadrat", args: struct {
			address string
		}{address: "N26"}, want: "", want1: false},
		{name: "mannheim-deluxe-plus-multiline", args: struct {
			address string
		}{address: "kommt alle nach mannheim , G3 4, Am Arsch GmbH"}, want: "g3", want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := detectQuadratAddress(tt.args.address)
			if got != tt.want {
				t.Errorf("detectQuadratAddress() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("detectQuadratAddress() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNormalizeAddressNewlineSplitter(t *testing.T) {
	de, _ := NewDE()
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "mannheim-deluxe", args: struct {
			address string
		}{address: "XyZ GmbH, Mannheim, G 3 4"}, want: []string{"g3", "xyz gmbh", "mannheim"}, wantErr: false},
		{name: "new-line", args: struct {
			address string
		}{address: "Eine GmbH \n Fürstenstr. 17 A"}, want: []string{"eine gmbh", "fuerstenstraße"}, wantErr: false},
		{name: "sankt", args: struct {
			address string
		}{address: "St.  Johan-Str. 111B"}, want: []string{"sankt johan straße"}, wantErr: false},
		{name: "strasze", args: struct {
			address string
		}{address: "Mönch-Strasze"}, want: []string{"moench straße"}, wantErr: false},
		{name: "strasse", args: struct {
			address string
		}{address: "Somestrasse 41  |  Im Tal der 100 Tränen"}, want: []string{"somestraße", "im tal der traenen"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := de.Street(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeAddressNewlineSplitter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeAddressNewlineSplitter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDE_Street(t *testing.T) {
	type fields struct {
		reCity        *regexp.Regexp
		rePostalCode  *regexp.Regexp
		reStreet      *regexp.Regexp
		reFoundStreet *regexp.Regexp
	}

	tests := []struct {
		name    string
		fields  fields
		address string
		want    []string
		wantErr bool
	}{
		{
			name:    "str. replacement",
			address: "Mühlstr. 34",
			want:    []string{"muehlstraße"},
			wantErr: false,
		},
		{
			name:    "strasse replacement",
			address: "Mühlstrasse 34",
			want:    []string{"muehlstraße"},
			wantErr: false,
		},
		{
			name:    "str replacement",
			address: "Mühlstr 34",
			want:    []string{"muehlstraße"},
			wantErr: false,
		},
		{
			name:    "str in the middle",
			address: "Industriestr 34",
			want:    []string{"industriestraße"},
			wantErr: false,
		},
		{
			name:    "str in the middle 2",
			address: "Zum Fingstried 34",
			want:    []string{"zum fingstried"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			de, _ := NewDE()
			got, err := de.Street(tt.address)
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
