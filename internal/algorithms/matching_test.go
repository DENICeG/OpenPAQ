// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package algorithms

import (
	"github.com/hbollon/go-edlib"
	"reflect"
	"testing"
)

func TestGetMatches(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		compareTo []string
		config    MatchSeverityConfig
		expected  []MatchResult
		shouldErr bool
	}{
		{
			name:      "Simple Match",
			input:     "word",
			compareTo: []string{"word"},
			config: MatchSeverityConfig{
				Algorithm:          edlib.Lcs,
				AlgorithmThreshold: 1,
			},
			expected: []MatchResult{
				{
					Value:      "word",
					Similarity: 1,
					WasPartial: false,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "No Match",
			input:     "stuff",
			compareTo: []string{"word"},
			config: MatchSeverityConfig{
				Algorithm:          edlib.Lcs,
				AlgorithmThreshold: 1,
			},
			expected:  nil,
			shouldErr: true,
		},
		// next test
		{
			name:      "No matches - threshold to high",
			input:     "wordz",
			compareTo: []string{"word"},
			config: MatchSeverityConfig{
				Algorithm:          edlib.Lcs,
				AlgorithmThreshold: 1,
			},
			expected:  nil,
			shouldErr: true,
		},
		// next test
		{
			name:      "Almost the right word",
			input:     "wordz",
			compareTo: []string{"word"},
			config: MatchSeverityConfig{
				Algorithm:          edlib.Lcs,
				AlgorithmThreshold: 0.8,
			},
			expected: []MatchResult{
				{
					Value:      "word",
					Similarity: 0.8,
					WasPartial: false,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Multiple response",
			input:     "word",
			compareTo: []string{"word", "word2", "word33"},
			config: MatchSeverityConfig{
				Algorithm:          edlib.Lcs,
				AlgorithmThreshold: 0.6,
			},
			expected: []MatchResult{
				{
					Value:      "word",
					Similarity: 1,
					WasPartial: false,
				},
				{
					Value:      "word2",
					Similarity: 0.8,
					WasPartial: false,
				},
				{
					Value:      "word33",
					Similarity: 0.6666667,
					WasPartial: false,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Simple Match - Allow partial match",
			input:     "word",
			compareTo: []string{"word"},
			config: MatchSeverityConfig{
				Algorithm:          edlib.Lcs,
				AlgorithmThreshold: 1,
				AllowPartialMatch:  true,
			},
			expected: []MatchResult{
				{
					Value:      "word",
					Similarity: 1,
					WasPartial: false,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Simple Match - Allow partial match",
			input:     "word asdf",
			compareTo: []string{"word", "word asdfg"},
			config: MatchSeverityConfig{
				Algorithm:                 edlib.Lcs,
				AlgorithmThreshold:        1,
				PartialAlgorithmThreshold: 1,
				AllowPartialMatch:         true,
			},
			expected: []MatchResult{
				{
					Value:      "word",
					Similarity: 0.44444445,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Multiple Match - Allow partial match",
			input:     "word asdf",
			compareTo: []string{"word", "word asdf"},
			config: MatchSeverityConfig{
				Algorithm:          edlib.Lcs,
				AlgorithmThreshold: 1,
				AllowPartialMatch:  true,
			},
			expected: []MatchResult{
				{
					Value:      "word asdf",
					Similarity: 1,
					WasPartial: false,
				},
				{
					Value:      "word",
					Similarity: 0.44444445,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Single Match - Allow partial match",
			input:     "word asdf",
			compareTo: []string{"word1"},
			config: MatchSeverityConfig{
				Algorithm:                 edlib.Lcs,
				AlgorithmThreshold:        1,
				AllowPartialMatch:         true,
				PartialAlgorithm:          edlib.Lcs,
				PartialAlgorithmThreshold: 0.5,
			},
			expected: []MatchResult{
				{
					Value:      "word1",
					Similarity: 0.33333334,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Partial Match - AllowPartialCompareListMatch",
			input:     "word",
			compareTo: []string{"stuff before word1", "word two"},
			config: MatchSeverityConfig{
				Algorithm:                    edlib.Lcs,
				AlgorithmThreshold:           0.9,
				PartialAlgorithm:             edlib.Lcs,
				PartialAlgorithmThreshold:    0.8,
				AllowPartialMatch:            true,
				AllowPartialCompareListMatch: true,
			},
			expected: []MatchResult{
				{
					Value:      "word two",
					Similarity: 0.5,
					WasPartial: true,
				},
				{
					Value:      "stuff before word1",
					Similarity: 0.22222222,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Partial Match-AllowPartialCompareListMatch-1",
			input:     "word",
			compareTo: []string{"stuff before word1", "word two"},
			config: MatchSeverityConfig{
				Algorithm:                    edlib.Lcs,
				AlgorithmThreshold:           0.9,
				PartialAlgorithm:             edlib.Lcs,
				PartialAlgorithmThreshold:    1,
				AllowPartialMatch:            true,
				AllowPartialCompareListMatch: true,
			},
			expected: []MatchResult{
				{
					Value:      "word two",
					Similarity: 0.5,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Partial Match-AllowPartialCompareListMatch-2",
			input:     "stuff word",
			compareTo: []string{"stuff before word1", "word two"},
			config: MatchSeverityConfig{
				Algorithm:                    edlib.Lcs,
				AlgorithmThreshold:           0.9,
				PartialAlgorithm:             edlib.Lcs,
				PartialAlgorithmThreshold:    0.8,
				AllowPartialMatch:            true,
				AllowPartialCompareListMatch: true,
			},
			expected: []MatchResult{
				{
					Value:      "stuff before word1",
					Similarity: 0.5555556,
					WasPartial: true,
				},
				{
					Value:      "word two",
					Similarity: 0,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Partial Match-AllowPartialCompareListMatch-3",
			input:     "stuff word",
			compareTo: []string{"stuff before word1", "word1 two"},
			config: MatchSeverityConfig{
				Algorithm:                    edlib.Lcs,
				AlgorithmThreshold:           0.9,
				PartialAlgorithm:             edlib.Lcs,
				PartialAlgorithmThreshold:    1,
				AllowPartialMatch:            true,
				AllowPartialCompareListMatch: true,
			},
			expected: []MatchResult{
				{
					Value:      "stuff before word1",
					Similarity: 0.5555556,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Partial Match-AllowPartialCompareListMatch-4",
			input:     "stuff word",
			compareTo: []string{"stuff before word1", "word two"},
			config: MatchSeverityConfig{
				Algorithm:                    edlib.Lcs,
				AlgorithmThreshold:           0.9,
				PartialAlgorithm:             edlib.Lcs,
				PartialAlgorithmThreshold:    1,
				AllowPartialMatch:            true,
				AllowPartialCompareListMatch: true,
				PartialExcludeWords:          []string{"stuff"},
			},
			expected: []MatchResult{
				{
					Value:      "word two",
					Similarity: 0,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
		// next test
		{
			name:      "Partial Match-AllowPartialCompareListMatch",
			input:     "stuff before word stuff behind",
			compareTo: []string{"word stuff"},
			config: MatchSeverityConfig{
				Algorithm:                          edlib.Lcs,
				AlgorithmThreshold:                 0.9,
				PartialAlgorithm:                   edlib.Lcs,
				PartialAlgorithmThreshold:          1,
				AllowPartialMatch:                  true,
				AllowCombineAllForwardCombinations: true,
			},
			expected: []MatchResult{
				{
					Value:      "word stuff",
					Similarity: 0.33333334,
					WasPartial: true,
				},
			},
			shouldErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, err := GetMatches(test.input, test.compareTo, test.config)

			if err != nil && !test.shouldErr {
				t.Errorf("got unexpected error: %s", err)
			}

			if !reflect.DeepEqual(test.expected, res) {
				t.Errorf("GetMatches() got: %v, expected %v", res, test.expected)
			}
		})
	}
}
