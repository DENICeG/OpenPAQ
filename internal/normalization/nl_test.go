// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCodeNL(t *testing.T) {
	nlNormalizer, _ := newNL()

	countryCode := nlNormalizer.GetCountryCode()

	if countryCode != "nl" {
		t.Errorf("got %s, want nl", countryCode)
	}
}

func TestCityNL(t *testing.T) {
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

func TestPostalCodeNoErrorNL(t *testing.T) {
	nlNormalizer, _ := newNL()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "whitespace in", input: "9865 MD", expect: "9865 MD"},
		{name: "whitespace should be in between", input: "9865MD", expect: "9865 MD"},
		{name: "special char remove", input: "9865MD+/(){}[]<>!§$%&=?*#€¿_\",:;-", expect: "9865 MD"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := nlNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestPostalCodeExpectErrorNL(t *testing.T) {
	nlNormalizer, _ := newNL()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Postal code lowercase", input: "6825mdd", expect: "6825MDD"},
		{name: "Postal code to long", input: "6825 MDD", expect: "6825 MDD"},
		{name: "Postal code to short", input: "825MD", expect: "825MD"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := nlNormalizer.PostalCode(test.input)
			if err == nil {
				t.Error("expect error got none")
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreetNL(t *testing.T) {
	esNormalizer, _ := newNL()

	tests := []struct {
		name            string
		input           string
		expectedContent []string
	}{
		{name: "no separator", input: "JustAStreet", expectedContent: []string{"justastreet"}},
		{name: "spaces", input: "Just A Street", expectedContent: []string{"just street"}},
		{name: "comma", input: "Just , Street", expectedContent: []string{"just", "street"}},
		{name: "new line", input: "Just \n Street", expectedContent: []string{"just", "street"}},
		{name: "remove special characters", input: "Just Street +/(){}[]<>!§'$%&=?*#€¿_\":;0123456789", expectedContent: []string{"just street"}},
		{name: "replace short word straat", input: "Just Street str.", expectedContent: []string{"just street straat"}},
		{name: "replace short word laan", input: "Just Street ln.", expectedContent: []string{"just street laan"}},
		{name: "replace short word weg", input: "Just Street wg.", expectedContent: []string{"just street weg"}},
		{name: "replace short word plein", input: "Just Street pl.", expectedContent: []string{"just street plein"}},
		{name: "replace short word gracht", input: "Just Street gr.", expectedContent: []string{"just street gracht"}},
		{name: "replace short word singel", input: "Just Street sgl.", expectedContent: []string{"just street singel"}},
		{name: "replace short word kade", input: "Just Street kd.", expectedContent: []string{"just street kade"}},
		{name: "replace short word hof", input: "Just Street hf.", expectedContent: []string{"just street hof"}},
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
