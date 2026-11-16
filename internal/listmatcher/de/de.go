// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package de

import (
	"openPAQ/internal/algorithms"
	"openPAQ/internal/normalization"
	"openPAQ/internal/types"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

type NormalizeEnclosure struct {
	Raw        string
	Normalized string
}

type CityPostalCodeItems struct {
	City        string
	PostalCodes []NormalizeEnclosure
}

type PostalCodeStreetItems struct {
	PostalCode string
	Streets    []NormalizeEnclosure
}

type DE struct {
	db            Database
	matcherConfig algorithms.MatchSeverityConfig

	postalCodeStreet  map[string]PostalCodeStreetItems
	cityPostalCode    map[string]CityPostalCodeItems
	firstLetterCities map[string][]string //key first letter of city

	normalize normalization.Normalize
}

func NewDE(database Database, normalizer normalization.Normalize, matcherConfig algorithms.MatchSeverityConfig) *DE {
	de := &DE{
		normalize:     normalizer,
		matcherConfig: matcherConfig,
		db:            database,
	}

	de.buildMatchingLists()

	return de
}

func (de *DE) buildMatchingLists() {

	logrus.Info("Building PostalCodeStreet for DE")
	postalCodeStreet, err := de.db.GetPostalCodeStreet(de.normalize)

	if err != nil {
		panic(err)
	}

	de.postalCodeStreet = postalCodeStreet

	logrus.Info("Building CityPostalCode for DE")

	cityPostalCode, err := de.db.GetCityPostalCode(de.normalize)
	if err != nil {
		panic(err)
	}

	de.cityPostalCode = cityPostalCode

	logrus.Info("Building FirstLetterCities for DE")

	de.buildFirstLetterCities()

	logrus.Info("Building indexes for DE finished")

}

func (de *DE) buildFirstLetterCities() {

	var cities []string

	for city := range de.cityPostalCode {
		cities = append(cities, city)
	}

	slices.Sort(cities)
	cities = slices.Compact(cities)

	de.firstLetterCities = make(map[string][]string)

	for _, city := range cities {

		cityFirstLetter := string([]rune(city)[0])

		if _, ok := de.firstLetterCities[cityFirstLetter]; !ok {
			de.firstLetterCities[cityFirstLetter] = []string{city}
			continue
		}

		de.firstLetterCities[cityFirstLetter] = append(de.firstLetterCities[cityFirstLetter], city)

	}
}

func (de *DE) GetCountryCode() string {
	return "de"
}

func (de *DE) CityStreetCheck(input types.NormalizeInput) chan types.PairMatching {

	reChan := make(chan types.PairMatching, 1)

	go func(c chan types.PairMatching) {
		defer close(c)
		result := types.PairMatching{}

		cityConfig := de.matcherConfig
		cityConfig.AlgorithmThreshold = 0.9
		cityConfig.AllowPartialMatch = true
		cityConfig.AllowPartialCompareListMatch = true
		cityConfig.PartialInputSeparators = []string{" ", "/", "-"}

		if input.City == "" {
			c <- result
			return
		}

		inputCityFirstLetter := string([]rune(input.City)[0])

		cityCandidates, err := algorithms.GetMatches(input.City, de.firstLetterCities[inputCityFirstLetter], cityConfig)
		if err != nil {
			c <- result
			return
		}

		var matches []types.CityStreetPostalCode

		for _, cityCandidate := range cityCandidates {
			cityPostalCodeItem, ok := de.cityPostalCode[cityCandidate.Value]

			if ok {
				for _, postalCode := range cityPostalCodeItem.PostalCodes {
					postalCodeStreetItems, foundPostalCode := de.postalCodeStreet[postalCode.Normalized]

					if foundPostalCode {

						res := de.PostalCodeStreetCheck(types.NormalizeInput{
							Streets:    input.Streets,
							PostalCode: postalCode.Normalized,
						})

						a := <-res

						if a.PostalCodeStreetMatch {
							for _, street := range a.PostalCodeStreetMatches {
								matches = append(matches, types.CityStreetPostalCode{
									City:                  cityPostalCodeItem.City,
									Street:                street.Street,
									PostalCode:            postalCodeStreetItems.PostalCode,
									CountryCode:           "de",
									StreetSimilarity:      street.StreetSimilarity,
									WasPartialStreetMatch: street.WasPartialStreetMatch,
									CitySimilarity:        cityCandidate.Similarity,
									WasPartialCityMatch:   cityCandidate.WasPartial,
									WasListMatch:          true,
								})
							}
						}
					}
				}
			}
		}

		if len(matches) > 0 {
			result = types.PairMatching{
				StreetCityMatch:   true,
				StreetCityMatches: matches,
			}
		}

		c <- result
	}(reChan)

	return reChan
}

func (de *DE) PostalCodeStreetCheck(input types.NormalizeInput) chan types.PairMatching {

	reChan := make(chan types.PairMatching, 1)

	go func(c chan types.PairMatching) {
		defer close(c)
		result := types.PairMatching{}

		postalCodeItem, ok := de.postalCodeStreet[input.PostalCode]

		var normalizedStreets []string

		for _, streetFromList := range postalCodeItem.Streets {
			normalizedStreets = append(normalizedStreets, strings.ReplaceAll(streetFromList.Normalized, "straße", ""))
		}

		var matches []types.PostalCodeStreet

		if ok {
			for _, inputStreetItem := range input.Streets {

				streetConfig := de.matcherConfig
				streetConfig.AllowPartialMatch = true
				streetConfig.AllowCombineAllForwardCombinations = true
				streetConfig.AlgorithmThreshold = streetConfig.DeListMatchAlgorithmThreshold
				if streetConfig.DeListMatchAlgorithmThreshold > streetConfig.PartialAlgorithmThreshold {
					streetConfig.PartialAlgorithmThreshold = streetConfig.DeListMatchAlgorithmThreshold
				}

				inputStreetItem = strings.ReplaceAll(inputStreetItem, "straße", "")

				streetCandidates, err := algorithms.GetMatches(inputStreetItem, normalizedStreets, streetConfig)

				if err != nil {
					continue
				}

				for _, streetCandidate := range streetCandidates {
					for _, street := range postalCodeItem.Streets {
						if strings.ReplaceAll(street.Normalized, "straße", "") == streetCandidate.Value {

							matches = append(matches, types.PostalCodeStreet{
								PostalCode:            postalCodeItem.PostalCode,
								Street:                street.Raw,
								CountryCode:           "de",
								StreetSimilarity:      streetCandidate.Similarity,
								WasPartialStreetMatch: streetCandidate.WasPartial,
								WasListMatch:          true,
							})
						}
					}
				}
			}
		}

		if len(matches) > 0 {
			result = types.PairMatching{
				PostalCodeStreetMatch:   true,
				PostalCodeStreetMatches: matches,
			}
		}

		c <- result
	}(reChan)

	return reChan
}

func (de *DE) PostalCodeCityCheck(input types.NormalizeInput) chan types.PairMatching {
	reChan := make(chan types.PairMatching, 1)

	go func(c chan types.PairMatching) {
		defer close(c)
		result := types.PairMatching{}

		cityConfig := de.matcherConfig
		cityConfig.AlgorithmThreshold = 0.9
		cityConfig.AllowPartialMatch = true
		cityConfig.AllowPartialCompareListMatch = true
		cityConfig.PartialInputSeparators = []string{" ", "/", "-"}

		if input.City == "" {
			c <- result
			return
		}

		inputCityFirstLetter := string([]rune(input.City)[0])
		cityCandidates, err := algorithms.GetMatches(input.City, de.firstLetterCities[inputCityFirstLetter], cityConfig)
		if err != nil {
			c <- result
			return
		}

		var matches []types.CityPostalCode

		for _, cityCandidate := range cityCandidates {
			cityPostalCodeItem := de.cityPostalCode[cityCandidate.Value]
			for _, postalCodeCandidate := range cityPostalCodeItem.PostalCodes {

				if postalCodeCandidate.Normalized == input.PostalCode {
					matches = append(matches, types.CityPostalCode{
						City:                cityPostalCodeItem.City,
						PostalCode:          postalCodeCandidate.Raw,
						CountryCode:         "de",
						CitySimilarity:      cityCandidate.Similarity,
						WasPartialCityMatch: cityCandidate.WasPartial,
						WasListMatch:        true,
					})
				}
			}
		}

		if len(matches) > 0 {
			result = types.PairMatching{
				CityPostalCodeMatch:   true,
				CityPostalCodeMatches: matches,
			}
		}

		c <- result

	}(reChan)

	return reChan
}
