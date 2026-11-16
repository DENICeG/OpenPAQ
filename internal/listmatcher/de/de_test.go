// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package de

import (
	"openPAQ/internal/algorithms"
	"openPAQ/internal/normalization"
	"openPAQ/internal/types"
	"reflect"
	"testing"

	"github.com/hbollon/go-edlib"
)

type mockDB struct {
	GetPostalCodeStreetResponse map[string]PostalCodeStreetItems
	GetCityPostalCodeResponse   map[string]CityPostalCodeItems
}

func (m mockDB) GetPostalCodeStreet(normalizer normalization.Normalize) (map[string]PostalCodeStreetItems, error) {
	return m.GetPostalCodeStreetResponse, nil

}

func (m mockDB) GetCityPostalCode(normalizer normalization.Normalize) (map[string]CityPostalCodeItems, error) {
	return m.GetCityPostalCodeResponse, nil
}

func getDbMock() mockDB {
	return mockDB{
		GetPostalCodeStreetResponse: map[string]PostalCodeStreetItems{
			"12345": {
				PostalCode: "12345",
				Streets: []NormalizeEnclosure{
					{
						Raw:        "Eins-Straße",
						Normalized: "eins straße",
					},
					{
						Raw:        "Zwei-Straße",
						Normalized: "zwei straße",
					},
					{
						Raw:        "Drei-Straße",
						Normalized: "drei straße",
					},
					{
						Raw:        "Eins-Gleich-Straße",
						Normalized: "eins gleich straße",
					},
					{
						Raw:        "EinsA-Gleich-Straße",
						Normalized: "einsa gleich straße",
					},
				},
			},
			"123456": {
				PostalCode: "123456",
				Streets: []NormalizeEnclosure{
					{
						Raw:        "Eins-Straße",
						Normalized: "eins straße",
					},
					{
						Raw:        "Zwei-Straße",
						Normalized: "zwei straße",
					},
				},
			},
			"52385": {
				PostalCode: "52385",
				Streets: []NormalizeEnclosure{
					{
						Raw:        "Eins-Straße",
						Normalized: "eins straße",
					},
					{
						Raw:        "Zwei-Straße",
						Normalized: "zwei straße",
					},
				},
			},
			"67433": {
				PostalCode: "67433",
				Streets: []NormalizeEnclosure{
					{
						Raw:        "Eins-Straße",
						Normalized: "eins straße",
					},
				},
			},
			"60311": {
				PostalCode: "60311",
				Streets: []NormalizeEnclosure{
					{
						Raw:        "Eins-Straße",
						Normalized: "eins straße",
					},
				},
			},
			"99999": {
				PostalCode: "99999",
				Streets: []NormalizeEnclosure{
					{
						Raw:        "Drei-Straße",
						Normalized: "drei straße",
					},
				},
			},
		},
		GetCityPostalCodeResponse: map[string]CityPostalCodeItems{
			"astadt": {
				City: "AStadt",
				PostalCodes: []NormalizeEnclosure{
					{
						Raw:        "12345",
						Normalized: "12345",
					},
				},
			},
			"bstadt": {
				City: "BStadt",
				PostalCodes: []NormalizeEnclosure{
					{
						Raw:        "123456",
						Normalized: "123456",
					},
				},
			},
			"nideggen": {
				City: "Nideggen",
				PostalCodes: []NormalizeEnclosure{
					{
						Raw:        "52385",
						Normalized: "52385",
					},
				},
			},
			"neustadt": {
				City: "Neustadt",
				PostalCodes: []NormalizeEnclosure{
					{
						Raw:        "67433",
						Normalized: "67433",
					},
				},
			},
			"frankfurt am main": {
				City: "Frankfurt am Main",
				PostalCodes: []NormalizeEnclosure{
					{
						Raw:        "60311",
						Normalized: "60311",
					},
				},
			},
			"frankfurt": {
				City: "Frankfurt",
				PostalCodes: []NormalizeEnclosure{
					{
						Raw:        "99999",
						Normalized: "99999",
					},
				},
			},
		},
	}

}

func TestCityStreetCheck(t *testing.T) {

	tests := []struct {
		name     string
		input    types.NormalizeInput
		expected types.PairMatching
	}{
		{
			name: "Full Match",
			input: types.NormalizeInput{
				Streets: []string{"eins straße"},
				City:    "astadt",
			},
			expected: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "AStadt",
						PostalCode:            "12345",
						Street:                "Eins-Straße",
						CountryCode:           "de",
						StreetSimilarity:      1,
						CitySimilarity:        1,
						WasPartialStreetMatch: false,
						WasListMatch:          true,
					},
				},
			},
		},
		// next test
		{
			name: "Street does not exist",
			input: types.NormalizeInput{
				Streets: []string{"letzte straße"},
				City:    "astadt",
			},
			expected: types.PairMatching{},
		},
		// next test
		{
			name: "No city passed",
			input: types.NormalizeInput{
				Streets: []string{"eins straße"},
				//City:    "astadt",
			},
			expected: types.PairMatching{},
		},
		// next test
		{
			name: "City does not exist",
			input: types.NormalizeInput{
				Streets: []string{"eins straße"},
				City:    "acityyyy",
			},
			expected: types.PairMatching{},
		},
		// next test
		{
			name: "Strange city name Match",
			input: types.NormalizeInput{
				Streets: []string{"eins straße"},
				City:    "nordrhein westfalen   nideggen  dueren",
			},
			expected: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                "Nideggen",
						Street:              "Eins-Straße",
						PostalCode:          "52385",
						CountryCode:         "de",
						StreetSimilarity:    1,
						CitySimilarity:      0.21052632,
						WasPartialCityMatch: true,
						WasListMatch:        true,
					},
				},
			},
		},
		// next test
		{
			name: "Strange city name Match 2",
			input: types.NormalizeInput{
				Streets: []string{"eins straße asdf"},
				City:    "neustadt a.d. weinstrasze",
			},
			expected: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "Neustadt",
						Street:                "Eins-Straße",
						PostalCode:            "67433",
						CountryCode:           "de",
						StreetSimilarity:      0.5,
						CitySimilarity:        0.32,
						WasPartialCityMatch:   true,
						WasPartialStreetMatch: true,
						WasListMatch:          true,
					},
				},
			},
		},
		// next test
		{
			name: "Two different cities with almost same name",
			input: types.NormalizeInput{
				Streets: []string{"eins straße"},
				City:    "frankfurt",
			},
			expected: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                "Frankfurt am Main",
						Street:              "Eins-Straße",
						PostalCode:          "60311",
						CountryCode:         "de",
						StreetSimilarity:    1,
						CitySimilarity:      0.5294118,
						WasPartialCityMatch: true,
						WasListMatch:        true,
					},
				},
			},
		},
	}

	deNormalizer, err := normalization.NewDE()
	if err != nil {
		t.Error("could not create normalizer for de")
	}

	// same config as used in implementation
	matcherConfig := algorithms.MatchSeverityConfig{
		Algorithm:                     edlib.Lcs,
		AlgorithmThreshold:            0.9,
		DeListMatchAlgorithmThreshold: 0.9,
		PartialAlgorithm:              edlib.Lcs,
		PartialAlgorithmThreshold:     1,
	}

	deListMatcher := NewDE(getDbMock(), deNormalizer, matcherConfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := <-deListMatcher.CityStreetCheck(tt.input)
			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("got %v, want %v", res, tt.expected)
			}
		})
	}
}

func TestPostalCodeStreetCheck(t *testing.T) {

	tests := []struct {
		name     string
		input    types.NormalizeInput
		expected types.PairMatching
	}{
		{
			name: "Full Match",
			input: types.NormalizeInput{
				Streets:    []string{"eins straße"},
				PostalCode: "12345",
			},
			expected: types.PairMatching{
				PostalCodeStreetMatch: true,
				PostalCodeStreetMatches: []types.PostalCodeStreet{
					{
						PostalCode:            "12345",
						Street:                "Eins-Straße",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
						WasListMatch:          true,
					},
				},
			},
		},
		// next test
		{
			name: "Multiple Match",
			input: types.NormalizeInput{
				Streets:    []string{"eins gleich straße"},
				PostalCode: "12345",
			},
			expected: types.PairMatching{
				PostalCodeStreetMatch: true,
				PostalCodeStreetMatches: []types.PostalCodeStreet{
					{
						PostalCode:            "12345",
						Street:                "Eins-Gleich-Straße",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
						WasListMatch:          true,
					},
					{
						PostalCode:            "12345",
						Street:                "EinsA-Gleich-Straße",
						CountryCode:           "de",
						StreetSimilarity:      0.9230769,
						WasPartialStreetMatch: false,
						WasListMatch:          true,
					},
				},
			},
		},
		// next test
		{
			name: "small typo in input",
			input: types.NormalizeInput{
				Streets:    []string{"eins gleicher straße"},
				PostalCode: "12345",
			},
			expected: types.PairMatching{
				PostalCodeStreetMatch: true,
				PostalCodeStreetMatches: []types.PostalCodeStreet{
					{
						PostalCode:            "12345",
						Street:                "Eins-Gleich-Straße",
						CountryCode:           "de",
						StreetSimilarity:      0.85714287,
						WasPartialStreetMatch: false,
						WasListMatch:          true,
					},
				},
			},
		},
		// next test
		{
			name: "Partial Match",
			input: types.NormalizeInput{
				Streets:    []string{"stuff before eins gleich straße stuff behind"},
				PostalCode: "12345",
			},
			expected: types.PairMatching{
				PostalCodeStreetMatch: true,
				PostalCodeStreetMatches: []types.PostalCodeStreet{
					{
						PostalCode:            "12345",
						Street:                "Eins-Gleich-Straße",
						CountryCode:           "de",
						StreetSimilarity:      0.31578946,
						WasPartialStreetMatch: true,
						WasListMatch:          true,
					},
				},
			},
		},
		// next test
		{
			name: "No Match",
			input: types.NormalizeInput{
				Streets:    []string{"unbekannt straße"},
				PostalCode: "postal code one",
			},
			expected: types.PairMatching{},
		},
	}

	deNormalizer, err := normalization.NewDE()
	if err != nil {
		t.Error("could not create normalizer for de")
	}

	// same config as used in implementation
	matcherConfig := algorithms.MatchSeverityConfig{
		Algorithm:                     edlib.Lcs,
		AlgorithmThreshold:            0.8,
		DeListMatchAlgorithmThreshold: 0.8,
		PartialAlgorithm:              edlib.Lcs,
		PartialAlgorithmThreshold:     1,
	}

	deListMatcher := NewDE(getDbMock(), deNormalizer, matcherConfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := <-deListMatcher.PostalCodeStreetCheck(tt.input)
			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("got %v, want %v", res, tt.expected)
			}
		})
	}
}

func TestBuildFirstLetterCities(t *testing.T) {

	de := DE{

		cityPostalCode: map[string]CityPostalCodeItems{
			"bcity":  {},
			"acity":  {},
			"bcity2": {},
			"ccity1": {},
			"ccity3": {},
			"ccity2": {},
		},
	}

	de.buildFirstLetterCities()

	expected := map[string][]string{
		"a": {
			"acity",
		},
		"b": {
			"bcity",
			"bcity2",
		},
		"c": {
			"ccity1",
			"ccity2",
			"ccity3",
		},
	}

	if !reflect.DeepEqual(de.firstLetterCities, expected) {
		t.Errorf("got %v, want %v", de.firstLetterCities, expected)
	}

}

func TestPostalCodeCityCheck(t *testing.T) {

	tests := []struct {
		name     string
		input    types.NormalizeInput
		expected types.PairMatching
	}{
		{
			name: "Full Match",
			input: types.NormalizeInput{
				City:       "astadt",
				PostalCode: "12345",
			},
			expected: types.PairMatching{
				CityPostalCodeMatch: true,
				CityPostalCodeMatches: []types.CityPostalCode{
					{
						City:                "AStadt",
						PostalCode:          "12345",
						CountryCode:         "de",
						CitySimilarity:      1,
						WasPartialCityMatch: false,
						WasListMatch:        true,
					},
				},
			},
		},
		// nest test
		{
			name: "No City provided",
			input: types.NormalizeInput{
				PostalCode: "12345",
			},
			expected: types.PairMatching{},
		},
		// nest test
		{
			name: "No PostalCode provided",
			input: types.NormalizeInput{
				City: "astadt",
			},
			expected: types.PairMatching{},
		},
		// nest test
		{
			name: "Unknown City",
			input: types.NormalizeInput{
				City:       "astadtzz",
				PostalCode: "12345",
			},
			expected: types.PairMatching{},
		},
		// nest test
		{
			name: "Wron PostalCode",
			input: types.NormalizeInput{
				City:       "astadt",
				PostalCode: "123456",
			},
			expected: types.PairMatching{},
		},
	}

	deNormalizer, err := normalization.NewDE()
	if err != nil {
		t.Error("could not create normalizer for de")
	}

	// same config as used in implementation
	matcherConfig := algorithms.MatchSeverityConfig{
		Algorithm:                 edlib.Lcs,
		AlgorithmThreshold:        0.9,
		PartialAlgorithm:          edlib.Lcs,
		PartialAlgorithmThreshold: 1,
	}

	deListMatcher := NewDE(getDbMock(), deNormalizer, matcherConfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := <-deListMatcher.PostalCodeCityCheck(tt.input)
			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("got %v, want %v", res, tt.expected)
			}
		})
	}
}

func TestGetCountryCode(t *testing.T) {

	de := DE{}

	res := de.GetCountryCode()

	if res != "de" {
		t.Errorf("got %v, want de", de.firstLetterCities)
	}

}
