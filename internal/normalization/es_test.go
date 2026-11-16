// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCodeES(t *testing.T) {
	esNormalizer, _ := newES()

	countryCode := esNormalizer.GetCountryCode()

	if countryCode != "es" {
		t.Errorf("got %s, want es", countryCode)
	}
}

func TestCityES(t *testing.T) {
	esNormalizer, _ := newES()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Want only small letters", input: "ACityWithCapitalLetters", expect: "acitywithcapitalletters"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := esNormalizer.City(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeNoErrorES(t *testing.T) {
	esNormalizer, _ := newES()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be removed", input: "12345abcdefghijklmnopqrstuvwxyz", expect: "12345"},
		{name: "Capital letters should be removed", input: "12345ABCDEFGHIJKLMNOPQRSTUVWXYZ", expect: "12345"},
		{name: "Special characters should be removed", input: "12345+/(){}[]<>!§$%&=?*#€¿_\",:;-", expect: "12345"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := esNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeExpectErrorES(t *testing.T) {
	esNormalizer, _ := newES()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Postal code to long", input: "123453", expect: "12345"},
		{name: "Postal code to short", input: "123", expect: "123"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := esNormalizer.PostalCode(test.input)
			if err == nil {
				t.Error("expect error got none")
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreetES(t *testing.T) {
	esNormalizer, _ := newES()

	tests := []struct {
		name            string
		input           string
		expectedContent []string
	}{
		{name: "no separator", input: "JustAStreet", expectedContent: []string{"justastreet"}},
		{name: "spaces", input: "Just A Street", expectedContent: []string{"just street"}},
		{name: "comma", input: "Just , Street", expectedContent: []string{"just", "street"}},
		{name: "new line", input: "Just \n Street", expectedContent: []string{"just", "street"}},
		{name: "c/o at the end", input: "Just Street c/o", expectedContent: []string{"just street"}},
		{name: "att at the end of word", input: "Just Streetatt", expectedContent: []string{"just streetatt"}},
		{name: "att at the end of word", input: "Just Streetatt", expectedContent: []string{"just streetatt"}},
		{name: "calle in street", input: "Calle Just", expectedContent: []string{"just"}},
		{name: "autopista in street", input: "autopista Just", expectedContent: []string{"just"}},
		{name: "special char in street", input: "cañada real just", expectedContent: []string{"just"}},
		{name: "Calle de la in street", input: "Calle de la Puebla de Farnals", expectedContent: []string{"puebla de farnals"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := esNormalizer.Street(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if !slices.Equal(test.expectedContent, res) {
				t.Errorf("got %v, want %v", res, test.expectedContent)
			}
		})
	}
}

func TestNormalizeStringES(t *testing.T) {
	esNormalizer, _ := newES()

	tests := []struct {
		name            string
		input           string
		expectedContent string
	}{
		{name: "remove spanish char ñ", input: "cañada real", expectedContent: "canada real"},
		{name: "remove spanish char ç", input: "caçada real", expectedContent: "cacada real"},
		{name: "remove spanish char ó", input: "caóada real", expectedContent: "caoada real"},
		{name: "remove spanish char ñ", input: "cañada real", expectedContent: "canada real"},
		{name: "remove spanish char ú", input: "caúada real", expectedContent: "cauada real"},
		{name: "remove spanish char á", input: "caáada real", expectedContent: "caaada real"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := esNormalizer.ReplaceSpanishLetters(test.input)
			if test.expectedContent != res {
				t.Errorf("got %s, want %s", res, test.expectedContent)
			}
		})
	}
}
