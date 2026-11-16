// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package normalization

import (
	"slices"
	"testing"
)

func TestGetCountryCodeUS(t *testing.T) {
	usNormalizer, _ := newUS()

	countryCode := usNormalizer.GetCountryCode()

	if countryCode != "us" {
		t.Errorf("got %s, want dk", countryCode)
	}
}

func TestPostalCodeUS(t *testing.T) {
	usNormalizer, _ := newUS()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be removed", input: "AA 12345asdf", expect: "12345"},
		{name: "Letters should be removed", input: "AA 12345 asdf", expect: "12345"},
		{name: "Numbers be should be removed after a dash", input: "acdsvsdv 12345-1234", expect: "12345"},
		{name: "special characters should be removed", input: "acdsvsdv 12345 !@#$%^&*()", expect: "12345"},
		{name: "special non numbers should be removed", input: " ASDF 12345", expect: "12345"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := usNormalizer.PostalCode(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestCityUS(t *testing.T) {
	usNormalizer, _ := newUS()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "Letters should be lower case", input: "Los Angeles", expect: "los angeles"},
		{name: "Whitespaces should be trimmed", input: "   Los Angeles   ", expect: "los angeles"},
		{name: "Special characters should be removed", input: "Los Ang+/(){}[]<>!§'$%&=?*#€¿_\":;ele12345s   ", expect: "los angeles"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := usNormalizer.City(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if res != test.expect {
				t.Errorf("got %s, want %s", res, test.expect)
			}
		})
	}
}

func TestStreetUS(t *testing.T) {
	usNormalizer, _ := newUS()

	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{name: "split at new line", input: "street \n part2", expect: []string{"street", "part"}},
		{name: "split at comma", input: "street , part2", expect: []string{"street", "part"}},
		{name: "replace north", input: "street n", expect: []string{"street north"}},
		{name: "replace south", input: "street s", expect: []string{"street south"}},
		{name: "replace west", input: "street w", expect: []string{"street west"}},
		{name: "replace east", input: "street e", expect: []string{"street east"}},
		{name: "replace northeast", input: "street ne", expect: []string{"street northeast"}},
		{name: "replace southeast", input: "street se", expect: []string{"street southeast"}},
		{name: "replace northwest", input: "street nw", expect: []string{"street northwest"}},
		{name: "replace southwest", input: "street sw", expect: []string{"street southwest"}},
		{name: "not replace e at the end of a word", input: "streete", expect: []string{"streete"}},
		{name: "delete floor and suite", input: "street floor suite", expect: []string{"street"}},
		{name: "replace st", input: "someSt st", expect: []string{"somest street"}},
		{name: "replace dr", input: "someDr dr", expect: []string{"somedr drive"}},
		{name: "replace rd", input: "somerd rd", expect: []string{"somerd road"}},
		{name: "replace rd", input: "someave ave", expect: []string{"someave avenue"}},
		{name: "replace blvd", input: "someblvd blvd", expect: []string{"someblvd boulevard"}},
		{name: "replace hwy", input: "somehwy hwy", expect: []string{"somehwy highway"}},
		{name: "replace pkwy", input: "somepkwy pkwy", expect: []string{"somepkwy parkway"}},
		{name: "replace number - test1", input: "8825 n. 23rd Ave Suite 100", expect: []string{"north 23rd avenue"}},
		{name: "replace number - test2", input: "1660 Whatever Ave", expect: []string{"whatever avenue"}},
		{name: "replace number - test3", input: "8825 n. 23rd Ave Suite 100 22nd rd asdf", expect: []string{"north 23rd avenue 22nd road asdf"}},
		{name: "replace number - test4", input: "477 Broadway 2 nd FL", expect: []string{"broadway"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := usNormalizer.Street(test.input)
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}

			if !slices.Equal(test.expect, res) {
				t.Errorf("got %v, want %v", res, test.expect)
			}
		})
	}
}
