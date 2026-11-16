// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package algorithms

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/hbollon/go-edlib"
)

type MatchResult struct {
	Value      string
	Similarity float32
	WasPartial bool
}

func GetMatches(input string, compareTo []string, config MatchSeverityConfig) ([]MatchResult, error) {
	var result []MatchResult

	matches, err := calculateSimilarity(input, compareTo, config.Algorithm, config.AlgorithmThreshold)
	if err == nil {
		result = append(result, matches...)
	}

	if config.AllowPartialMatch {
		matches, err = getPartialMatches(input, compareTo, config)
		if err == nil {
			result = append(result, matches...)
		}
	}

	if len(result) > 0 {
		cleanedResults := cleanMatchResults(result)
		return cleanedResults, nil
	}

	return result, errors.New("no matches found")
}

func getPartialMatches(input string, compareTo []string, config MatchSeverityConfig) ([]MatchResult, error) {
	var result []MatchResult

	var splitInputs []string

	if len(config.PartialInputSeparators) == 0 {
		config.PartialInputSeparators = []string{" "}
	}

	if config.AllowCombineAllForwardCombinations {
		splitInputs = combineAllForwardCombinations(input, config.PartialInputSeparators)
	} else {
		splitInputs = splitString(input, config.PartialInputSeparators)
	}

	splitInputs = filterList(splitInputs, config.PartialExcludeWords)
	compareTo = filterList(compareTo, config.PartialExcludeWords)

	for _, splitInput := range splitInputs {

		res, _ := calculateSimilarity(splitInput, compareTo, config.PartialAlgorithm, config.PartialAlgorithmThreshold)

		if len(res) > 0 {
			for _, item := range res {

				var err error
				item.Similarity, err = edlib.StringsSimilarity(input, item.Value, config.PartialAlgorithm)
				if err == nil {
					item.WasPartial = true
					result = append(result, item)
				}
			}
		}
	}

	if config.AllowPartialCompareListMatch {

		var splitCompareTo []string
		if len(config.PartialCompareListSeparators) == 0 {
			config.PartialCompareListSeparators = []string{" "}
		}

		for _, splitInput := range splitInputs {
			for _, compareItem := range compareTo {
				splitCompareTo = splitString(compareItem, config.PartialCompareListSeparators)
				splitCompareTo = filterList(splitCompareTo, config.PartialExcludeWords)

				res, _ := calculateSimilarity(splitInput, splitCompareTo, config.PartialAlgorithm, config.PartialAlgorithmThreshold)

				if len(res) > 0 {
					for _, item := range res {

						var err error
						item.Similarity, err = edlib.StringsSimilarity(input, compareItem, config.PartialAlgorithm)
						if err == nil {
							item.WasPartial = true
							item.Value = compareItem
							result = append(result, item)
						}
					}
				}

			}
		}

	}

	if len(result) > 0 {
		return result, nil
	}

	return nil, fmt.Errorf("no match found")

}

func filterList(input, filter []string) []string {
	var result []string

	for _, candidate := range input {
		if slices.Contains(filter, candidate) {
			continue
		}
		result = append(result, candidate)

	}

	return result
}

func calculateSimilarity(valueToTest string, allowedValues []string, algorithm edlib.Algorithm, threshold float32) ([]MatchResult, error) {
	res, err := edlib.FuzzySearchSet(valueToTest, allowedValues, 10, algorithm)

	if err != nil || res == nil {
		return nil, errors.New("no match found")
	}

	var cleanedList []string

	for _, item := range res {
		if item != "" {
			cleanedList = append(cleanedList, item)
		}
	}

	var result []MatchResult

	for _, candidate := range cleanedList {

		similarity, err := edlib.StringsSimilarity(valueToTest, candidate, algorithm)

		if err != nil {
			continue
		}

		if similarity >= threshold {

			result = append(result, MatchResult{
				Value:      candidate,
				Similarity: similarity,
			})
		}

	}

	return result, nil

}

func cleanMatchResults(input []MatchResult) []MatchResult {
	var cleanedResults []MatchResult
	knownResults := make(map[string]MatchResult)

	for _, res := range input {

		knownResult, ok := knownResults[res.Value]

		if !ok {
			knownResults[res.Value] = res
		}

		if knownResult.Similarity < res.Similarity {
			knownResults[res.Value] = res
		}
	}

	for _, res := range knownResults {
		cleanedResults = append(cleanedResults, res)
	}

	sort.Slice(cleanedResults, func(i, j int) bool {
		return cleanedResults[i].Similarity > cleanedResults[j].Similarity
	})

	return cleanedResults
}

func combineAllForwardCombinations(input string, separators []string) []string {

	var listOfSingleWords []string
	for _, sep := range separators {

		listOfSingleWords = append(listOfSingleWords, strings.Split(input, sep)...)
	}

	result := make([]string, 0)

	for i := 0; i < len(listOfSingleWords); i++ {
		var temp []string
		for j, word := range listOfSingleWords {
			if i <= j {
				temp = append(temp, word)
				result = append(result, strings.Join(temp, " "))
			}
		}
	}

	return result
}

func splitString(input string, separators []string) []string {

	var result []string
	for _, sep := range separators {
		result = append(result, strings.Split(input, sep)...)
	}

	return result
}

type MatchSeverityConfig struct {
	Algorithm                          edlib.Algorithm `json:"matching_algorithm"`
	AlgorithmThreshold                 float32         `json:"matching_threshold"`
	AllowPartialMatch                  bool
	AllowPartialCompareListMatch       bool
	PartialAlgorithm                   edlib.Algorithm `json:"partial_matching_algorithm"`
	PartialAlgorithmThreshold          float32         `json:"partial_matching_threshold"`
	DeListMatchAlgorithmThreshold      float32         `json:"de_list_matching_algorithm_threshold"`
	PartialInputSeparators             []string
	PartialExcludeWords                []string
	PartialCompareListSeparators       []string
	AllowCombineAllForwardCombinations bool
	AllowedAmountOfChangedChars        int `json:"allowed_amount_of_changed_chars"`
}
