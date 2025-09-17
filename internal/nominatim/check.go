package nominatim

import (
	"context"
	"maps"
	"openPAQ/internal/algorithms"
	"openPAQ/internal/types"
	"strings"
	"sync"
)

func getInputLanguages(nominatimLanguages []string, inputCC string) []string {
	var additionalCountryCodes []string
	for _, i := range WorldWideLanguages {
		if i.cc == strings.ToUpper(inputCC) {
			additionalCountryCodes = i.languages
			break
		}
	}

	if len(additionalCountryCodes) != 0 {
		nominatimLanguages = append(nominatimLanguages, additionalCountryCodes...)
	}

	var lang = make(map[string]bool)

	for _, v := range nominatimLanguages {
		lang[v] = true
	}

	res := []string{}
	for v := range maps.Keys(lang) {
		res = append(res, v)
	}

	//nominatimLanguages = res

	return res
}

type nominatimRequestParameter struct {
	city        string
	street      string
	countryCode string
	postalCode  string
}

func (nom *Nominatim) request(ctx context.Context, reqParams nominatimRequestParameter) []ParsedResult {
	languages := getInputLanguages(nom.languages, reqParams.countryCode)

	writeProtection := sync.Mutex{}
	var parsedResults []ParsedResult

	wg := sync.WaitGroup{}

	for _, language := range languages {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()

			var results []NominatimCoreResult

			results, errSearchString := nom.api.RequestWithSearchString(ctx, nom.url, generateSearchString(reqParams.postalCode, reqParams.street, reqParams.city), "1", l)
			moreResults, errParameter := nom.api.RequestWithParameters(ctx, nom.url, NominatimDetailRequest{
				Street:     reqParams.street,
				PostalCode: reqParams.postalCode,
				City:       reqParams.city,
			}, "", l)

			if errSearchString != nil && errParameter != nil {
				return
			}

			results = append(results, moreResults...)

			for _, res := range results {
				normalizer := nom.normalizer.Get(strings.ToLower(res.CountryCode))
				parsedResult := res.parse()

				normalizedResult, err := parsedResult.normalize(normalizer)
				if err != nil {
					continue
				}

				writeProtection.Lock()
				parsedResults = append(parsedResults, normalizedResult)
				writeProtection.Unlock()
			}

		}(language)
	}

	wg.Wait()

	return parsedResults

}

func (nom *Nominatim) CityStreetCheck(ctx context.Context, input types.NormalizeInput) chan types.PairMatching {
	reChan := make(chan types.PairMatching, 1)
	go func(c chan types.PairMatching) {
		defer close(reChan)
		result := types.PairMatching{}

		var matches []types.CityStreetPostalCode

		for _, street := range input.Streets {

			nominatimResults := nom.request(ctx, nominatimRequestParameter{
				city:        input.City,
				street:      street,
				countryCode: input.CountryCode,
			})

			nominatimResults = removeDuplicatesParsedResult(nominatimResults)

			for _, nominatimResult := range nominatimResults {

				validateConfigStreet := nom.config

				if input.CountryCode == "pl" ||
					input.CountryCode == "es" ||
					input.CountryCode == "it" ||
					input.CountryCode == "de" ||
					input.CountryCode == "gb" {
					validateConfigStreet.AllowPartialMatch = true
					validateConfigStreet.AllowPartialCompareListMatch = true
				}

				streetMatch, streetMatchRrr := algorithms.GetMatches(street, nominatimResult.Street, validateConfigStreet)
				if streetMatchRrr != nil || len(streetMatch) < 1 {
					continue
				}

				validateConfigCity := nom.config
				validateConfigCity.AllowPartialMatch = true
				validateConfigCity.AllowPartialCompareListMatch = true

				cityMatch, cityMatchErr := algorithms.GetMatches(input.City, nominatimResult.City, validateConfigCity)

				if cityMatchErr != nil || len(cityMatch) < 1 {
					continue
				}

				for _, streetSep := range streetMatch {
					for _, city := range cityMatch {

						if city.Similarity > 0 && streetSep.Similarity > 0 {
							matches = append(matches, types.CityStreetPostalCode{
								City:                  city.Value,
								PostalCode:            nominatimResult.PostalCode,
								Street:                streetSep.Value,
								CountryCode:           nominatimResult.CountryCode,
								StreetSimilarity:      streetSep.Similarity,
								WasPartialStreetMatch: streetSep.WasPartial,
								CitySimilarity:        city.Similarity,
								WasPartialCityMatch:   city.WasPartial,
							})
						}
					}
				}
			}

		}

		uniqueMates := types.RemoveDuplicate(matches)

		if len(uniqueMates) > 0 {
			result.StreetCityMatch = true
			result.StreetCityMatches = append(result.StreetCityMatches, uniqueMates...)
		}
		c <- result
	}(reChan)
	return reChan
}

func (nom *Nominatim) PostalCodeStreetCheck(ctx context.Context, input types.NormalizeInput, cityStreet []types.CityStreetPostalCode) chan types.PairMatching {
	reChan := make(chan types.PairMatching, 1)
	go func(c chan types.PairMatching) {
		defer close(reChan)
		result := types.PairMatching{}

		for _, element := range cityStreet {

			if (strings.Contains(input.PostalCode, element.PostalCode) && element.PostalCode != "" || strings.Contains(element.PostalCode, input.PostalCode)) && input.PostalCode != "" {
				result.PostalCodeStreetMatch = true
				result.PostalCodeStreetMatches = append(result.PostalCodeStreetMatches, types.PostalCodeStreet{
					PostalCode:            element.PostalCode,
					Street:                element.Street,
					CountryCode:           element.CountryCode,
					StreetSimilarity:      element.StreetSimilarity,
					WasPartialStreetMatch: element.WasPartialStreetMatch,
				})
			}
		}

		if result.PostalCodeStreetMatch {
			c <- result
			return
		}

		var matches []types.PostalCodeStreet

		for _, street := range input.Streets {
			req := nominatimRequestParameter{
				street:      street,
				postalCode:  input.PostalCode,
				countryCode: input.CountryCode,
			}

			if input.CountryCode == "gb" {
				req = nominatimRequestParameter{
					street:      street,
					countryCode: input.CountryCode,
				}
			}

			nominatimResults := nom.request(ctx, req)

			nominatimResults = removeDuplicatesParsedResult(nominatimResults)

			validateConfigStreet := nom.config

			if input.CountryCode == "pl" ||
				input.CountryCode == "es" ||
				input.CountryCode == "it" ||
				input.CountryCode == "de" ||
				input.CountryCode == "gb" {
				validateConfigStreet.AllowPartialMatch = true
				validateConfigStreet.AllowPartialCompareListMatch = true
			}

			for _, nominatimResult := range nominatimResults {

				streetMatch, err := algorithms.GetMatches(street, nominatimResult.Street, validateConfigStreet)
				if err != nil {
					continue
				}

				if (strings.Contains(nominatimResult.PostalCode, input.PostalCode) && input.PostalCode != "" || strings.Contains(input.PostalCode, nominatimResult.PostalCode)) && nominatimResult.PostalCode != "" {

					for _, streetPart := range streetMatch {

						matches = append(matches, types.PostalCodeStreet{
							PostalCode:            nominatimResult.PostalCode,
							Street:                streetPart.Value,
							CountryCode:           nominatimResult.CountryCode,
							StreetSimilarity:      streetPart.Similarity,
							WasPartialStreetMatch: streetPart.WasPartial,
						})
					}
				}
			}

		}

		uniqueMates := types.RemoveDuplicate(matches)

		if len(matches) > 0 {
			result.PostalCodeStreetMatch = true
			result.PostalCodeStreetMatches = append(result.PostalCodeStreetMatches, uniqueMates...)
		}

		c <- result
	}(reChan)

	return reChan
}

func (nom *Nominatim) PostalCodeCityCheck(ctx context.Context, input types.NormalizeInput, cityStreet []types.CityStreetPostalCode) chan types.PairMatching {
	reChan := make(chan types.PairMatching, 1)

	go func(c chan types.PairMatching) {
		defer close(reChan)
		result := types.PairMatching{}

		for _, element := range cityStreet {

			if (strings.Contains(input.PostalCode, element.PostalCode) && element.PostalCode != "" || strings.Contains(element.PostalCode, input.PostalCode)) && input.PostalCode != "" {
				result.CityPostalCodeMatch = true
				result.CityPostalCodeMatches = append(result.CityPostalCodeMatches, types.CityPostalCode{
					PostalCode:          element.PostalCode,
					City:                element.City,
					CountryCode:         element.CountryCode,
					CitySimilarity:      element.CitySimilarity,
					WasPartialCityMatch: element.WasPartialCityMatch,
				})
			}
		}
		if result.CityPostalCodeMatch {
			c <- result
			return
		}

		var matches []types.CityPostalCode

		nominatimResults := nom.request(ctx, nominatimRequestParameter{
			postalCode:  input.PostalCode,
			countryCode: input.CountryCode,
		})

		moreNominatimResults := nom.request(ctx, nominatimRequestParameter{
			postalCode:  input.PostalCode,
			city:        input.City,
			countryCode: input.CountryCode,
		})

		nominatimResults = append(nominatimResults, moreNominatimResults...)

		nominatimResults = removeDuplicatesParsedResult(nominatimResults)

		validateConfigCity := nom.config
		validateConfigCity.AllowPartialMatch = true
		validateConfigCity.AllowPartialCompareListMatch = true

		for _, nominatimResult := range nominatimResults {

			cityMatch, err := algorithms.GetMatches(input.City, nominatimResult.City, validateConfigCity)
			if err != nil {
				continue
			}

			if (strings.Contains(nominatimResult.PostalCode, input.PostalCode) && input.PostalCode != "" || strings.Contains(input.PostalCode, nominatimResult.PostalCode)) && nominatimResult.PostalCode != "" {
				for _, city := range cityMatch {

					matches = append(matches, types.CityPostalCode{
						City:                city.Value,
						PostalCode:          nominatimResult.PostalCode,
						CountryCode:         nominatimResult.CountryCode,
						CitySimilarity:      city.Similarity,
						WasPartialCityMatch: city.WasPartial,
					})
				}
			}
		}

		uniqueMates := types.RemoveDuplicate(matches)

		if len(uniqueMates) > 0 {
			result.CityPostalCodeMatch = true
			result.CityPostalCodeMatches = append(result.CityPostalCodeMatches, uniqueMates...)
		}
		c <- result
	}(reChan)
	return reChan
}

func generateSearchString(args ...string) string {
	return strings.Join(args, ",")
}

func removeDuplicatesParsedResult(nomRes []ParsedResult) []ParsedResult {

	temp := make(map[string]ParsedResult)
	for _, res := range nomRes {
		temp[res.hash()] = res
	}

	var results []ParsedResult

	for _, res := range temp {
		results = append(results, res)
	}
	return results

}
