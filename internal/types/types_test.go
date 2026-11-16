// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package types

import (
	"openPAQ/internal/normalization"
	"reflect"
	"testing"
)

func TestRemoveDuplicate(t *testing.T) {

	inputComparables := []CityPostalCode{{
		City:        "StadtA",
		PostalCode:  "12345",
		CountryCode: "de",
	},
		{
			City:        "StadtB",
			PostalCode:  "11111",
			CountryCode: "de",
		}, {
			City:        "StadtA",
			PostalCode:  "12356",
			CountryCode: "de",
		}}

	inputStrings := []string{"a", "c", "z", "y", "w"}

	t.Run("duplicate-structures", func(t *testing.T) {
		duplicates := append(inputComparables, inputComparables...)
		got := RemoveDuplicate(duplicates)
		want := inputComparables
		if !reflect.DeepEqual(got, want) {
			t.Errorf("match() got = %v, want %v", got, want)
		}
	})

	t.Run("multi-duplicate-strings", func(t *testing.T) {
		duplicates := append(inputStrings, inputStrings...)
		duplicates = append(duplicates, duplicates...)
		got := RemoveDuplicate(duplicates)
		want := inputStrings
		if !reflect.DeepEqual(got, want) {
			t.Errorf("match() got = %v, want %v", got, want)
		}
	})

}

func TestInput_Normalize(t *testing.T) {
	deNormalizer, err := normalization.NewDE()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name     string
		rawInput Input
		want     NormalizeInput
	}{
		{name: "postal-code-fill", rawInput: Input{
			Street:      "Blubstraße",
			City:        "Köln",
			PostalCode:  "DE-1234",
			CountryCode: "DE",
			Normalizer:  deNormalizer,
		}, want: NormalizeInput{
			Streets:     []string{"blubstraße"},
			City:        "koeln",
			PostalCode:  "01234",
			CountryCode: "de",
		}},
		{name: "city", rawInput: Input{
			Street:      "Blubstraße",
			City:        "Hamburg Nord",
			PostalCode:  "DE44444",
			CountryCode: "DE",
			Normalizer:  deNormalizer,
		}, want: NormalizeInput{
			Streets:     []string{"blubstraße"},
			City:        "hamburg nord",
			PostalCode:  "44444",
			CountryCode: "de",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.rawInput
			if got := i.Normalize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInput_Normalize4Nominatim(t *testing.T) {
	tests := []struct {
		name     string
		rawInput Input
		want     Input
	}{{name: "all-inputs-lower-case", rawInput: Input{
		Street:      "A Street",
		City:        "A City",
		PostalCode:  "DE-56765",
		CountryCode: "DE",
	}, want: Input{
		Street:      "A Street",
		City:        "a city",
		PostalCode:  "de-56765",
		CountryCode: "de",
	}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.rawInput
			if got := i.Normalize4Nominatim(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize4Nominatim() = %v, want %v", got, tt.want)
			}
		})
	}
}
