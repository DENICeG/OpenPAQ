// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCodeAT(t *testing.T) {
	atNormalizer, _ := newAT()

	countryCode := atNormalizer.GetCountryCode()

	if countryCode != "at" {
		t.Errorf("got %s, want at", countryCode)
	}
}

func TestCityAT(t *testing.T) {
	atNormalizer, _ := newAT()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Want only small letters", input: "ACityWithCapitalLetters", expect: "acitywithcapitalletters"},
		{name: "Remove Sankt", input: "ACityWithCapitalLetters Sankt", expect: "acitywithcapitalletters "},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := atNormalizer.City(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

// please mach weiter hier
func TestPostalCodeNoErrorAT(t *testing.T) {
	atNormalizer, _ := newAT()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be removed", input: "1234abcdefghijklmnopqrstuvwxyz", expect: "1234"},
		{name: "Capital letters should be removed", input: "1234ABCDEFGHIJKLMNOPQRSTUVWXYZ", expect: "1234"},
		{name: "Special characters should be removed", input: "1234+/(){}[]<>!§$%&=?*#€¿_\",:;-", expect: "1234"},
		{name: "Strip Leading 0", input: "0234", expect: "234"},
		{name: "Don´t Strip anything", input: "1234", expect: "1234"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := atNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeExpectErrorAT(t *testing.T) {
	atNormalizer, _ := newAT()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Postal code to long", input: "12345", expect: "12345"},
		{name: "Postal code to short", input: "123", expect: "123"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := atNormalizer.PostalCode(test.input)
			if err == nil {
				t.Error("expect error got none")
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreetAT(t *testing.T) {
	atNormalizer, _ := newAT()

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
		{name: "Umlaute", input: "Jüst , Straße", expectedContent: []string{"juest", "strasse"}},
		{name: "Objektzusätze entfernen", input: "Just , Street, Tür 09", expectedContent: []string{"just", "street"}},
		{name: "Hausnummer Entfernen", input: "Just , Street 33a", expectedContent: []string{"just", "street"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := atNormalizer.Street(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if !slices.Equal(test.expectedContent, res) {
				t.Errorf("got %v, want %v", res, test.expectedContent)
			}
		})
	}
}
