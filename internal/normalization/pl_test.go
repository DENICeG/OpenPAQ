// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCodePl(t *testing.T) {
	usNormalizer, _ := newUS()

	countryCode := usNormalizer.GetCountryCode()

	if countryCode != "us" {
		t.Errorf("got %s, want dk", countryCode)
	}
}

func TestPostalCodePl(t *testing.T) {
	plNormalizer, _ := NewPl()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be removed", input: "AA 12345asdf", expect: "12-345"},
		{name: "five digits should be reformatted", input: "12345", expect: "12-345"},
		{name: "six digits should not be reformatted", input: "123456", expect: "123456"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := plNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestCityPl(t *testing.T) {
	plNormalizer, _ := NewPl()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be lower case", input: "IamBig", expect: "iambig"},
		{name: "Whitespaces should be trimmed", input: "   ALotOfSpaces   ", expect: "alotofspaces"},
		{name: "Special characters should be removed", input: "blub+/(){}[]<>!§'$%&=?*#€¿_\":;12345ber   ", expect: "blubber"},
		{name: "Special letters should be removed", input: "ęóąłżźść", expect: "eoalzzsc"},
		{name: "unicodes should be replaced", input: "Powsta\u0144cow \u015al\u0105skich 108/63", expect: "powstańcow slaskich"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := plNormalizer.City(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreetPl(t *testing.T) {
	plNormalizer, _ := NewPl()

	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{name: "Letters should be lower case", input: "IamBig", expect: []string{"iambig"}},
		{name: "Split at new line", input: "first \nsecond", expect: []string{"first", "second"}},
		{name: "Split at comma", input: "first,second", expect: []string{"first", "second"}},
		{name: "Special characters should be removed", input: "blub+/(){}[]<>!§'$%&=?*#€¿_\":;12345ber   ", expect: []string{"blubber"}},
		{name: "Special letters should be removed", input: "ęóąłżźść", expect: []string{"eoalzzsc"}},
		{name: "Remove words with dot at the end", input: "Ab. ft. bla blub", expect: []string{"bla blub"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := plNormalizer.Street(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if !slices.Equal(test.expect, res) {
				t.Errorf("got %v, want %v", res, test.expect)
			}
		})
	}
}
