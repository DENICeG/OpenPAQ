// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package listmatcher

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"openPAQ/internal/algorithms"
	"openPAQ/internal/listmatcher/de"
	types2 "openPAQ/internal/listmatcher/types"
	"openPAQ/internal/normalization"
	"openPAQ/internal/types"
)

type CountryMatcher interface {
	CityStreetCheck(input types.NormalizeInput) chan types.PairMatching
	PostalCodeStreetCheck(types.NormalizeInput) chan types.PairMatching
	PostalCodeCityCheck(types.NormalizeInput) chan types.PairMatching
	GetCountryCode() string
}

type ListMatcher struct {
	matcherConfig algorithms.MatchSeverityConfig
	checker       []CountryMatcher
}

func NewMatcher(matcherConfig algorithms.MatchSeverityConfig) *ListMatcher {
	diy := &ListMatcher{
		matcherConfig: matcherConfig,
	}

	return diy
}

func (lm *ListMatcher) Register(cc string, dbConfig types2.DatabaseConfig, matcherConfig algorithms.MatchSeverityConfig) error {

	if cc == "de" {
		deNormalizer, err := normalization.NewDE()
		if err != nil {
			panic(err)
		}

		lm.checker = append(lm.checker, de.NewDE(de.NewDatabase(dbConfig), deNormalizer, matcherConfig))
	}

	return nil
}

func (lm *ListMatcher) GetConfig() algorithms.MatchSeverityConfig {
	return lm.matcherConfig
}

func (lm *ListMatcher) getChecker(cc string) (CountryMatcher, error) {
	for _, i := range lm.checker {
		if i.GetCountryCode() == cc {
			return i, nil
		}
	}
	return nil, fmt.Errorf("unable to find matcher")
}

func (lm *ListMatcher) Possible(cc string) bool {
	if _, err := lm.getChecker(cc); err != nil {
		return false
	}
	return true
}

func (lm *ListMatcher) Handle(ctx context.Context, input types.Input) <-chan types.PairMatching {
	c := make(chan types.PairMatching, 1)
	go func() {
		defer close(c)
		var result types.PairMatching

		check, err := lm.getChecker(input.CountryCode)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err, "countryCode": input.CountryCode}).Error("unable to find country checker")
			c <- result
			return
		}

		normalizedInput := input.Normalize()
		rC := check.CityStreetCheck(normalizedInput)
		pC := check.PostalCodeCityCheck(normalizedInput)
		qC := check.PostalCodeStreetCheck(normalizedInput)

		for counter := 0; counter < 3; counter++ {
			select {
			case r := <-rC:
				rC = nil
				result.StreetCityMatch = r.StreetCityMatch
				result.StreetCityMatches = r.StreetCityMatches
			case p := <-pC:
				pC = nil
				result.CityPostalCodeMatch = p.CityPostalCodeMatch
				result.CityPostalCodeMatches = p.CityPostalCodeMatches
			case q := <-qC:
				qC = nil
				result.PostalCodeStreetMatch = q.PostalCodeStreetMatch
				result.PostalCodeStreetMatches = q.PostalCodeStreetMatches
			case <-ctx.Done():
				c <- result
				return
			}
		}
		c <- result
	}()
	return c
}
