package types

import (
	"openPAQ/internal/algorithms"
	"openPAQ/internal/normalization"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var InputNormalizerErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace:   "datascience",
	Subsystem:   "",
	Name:        "inputNormalizerErrorCounter",
	Help:        "",
	ConstLabels: nil,
}, []string{"type", "country_code"})

type Input struct {
	Street      string                  `json:"street"`
	City        string                  `json:"city"`
	PostalCode  string                  `json:"postal_code"`
	CountryCode string                  `json:"country_code"`
	Normalizer  normalization.Normalize `json:"-"`
}

type NormalizeInput struct {
	Streets     []string
	City        string
	PostalCode  string
	CountryCode string
}

func (i Input) Normalize() NormalizeInput {
	var res NormalizeInput

	res.CountryCode = strings.ReplaceAll(strings.ToLower(i.CountryCode), " ", "")

	newPlz, err := i.Normalizer.PostalCode(i.PostalCode)
	if err != nil {
		InputNormalizerErrorCounter.WithLabelValues("postal_code", res.CountryCode).Inc()
	}
	res.PostalCode = newPlz

	newCity, err := i.Normalizer.City(i.City)
	res.City = strings.ToLower(newCity)
	if err != nil || len(res.City) == 0 {
		InputNormalizerErrorCounter.WithLabelValues("city", res.City).Inc()
	}

	streetLower := strings.ToLower(i.Street)
	res.Streets, err = i.Normalizer.Street(streetLower)
	if err != nil || len(res.Streets) == 0 {
		InputNormalizerErrorCounter.WithLabelValues("street", i.Street).Inc()
	}

	return res
}

func (i Input) Normalize4Nominatim() Input {
	var res Input
	res.PostalCode = strings.ReplaceAll(strings.ToLower(i.PostalCode), " ", "")
	res.City = strings.ToLower(i.City)
	res.CountryCode = strings.ReplaceAll(strings.ToLower(i.CountryCode), " ", "")
	res.Street = i.Street
	res.Normalizer = i.Normalizer
	return res
}

type SourceOfTruth struct {
	StreetMatched           bool `json:"street_matched"`
	CityMatched             bool `json:"city_matched"`
	PostalCodeMatched       bool `json:"postal_code_matched"`
	CityToPostalCodeMatched bool `json:"city_to_postal_code_matched"`
	CountryCodeMatched      bool `json:"country_code_matched"`
}

type Result struct {
	Input
	SourceOfTruth
	Version string `json:"version"`
}

type DebugResult struct {
	Result
	DebugDetails `json:"details"`
}

type DebugDetails struct {
	algorithms.MatchSeverityConfig `json:"parameters"`
	StreetCityMatches              []CityStreetPostalCode `json:"city_street_matches"`
	StreetPostalCodeMatches        []PostalCodeStreet     `json:"postal_code_street_matches"`
	PostalCodeCityMatches          []CityPostalCode       `json:"city_postal_code_matches"`
}

type PairMatching struct {
	PostalCodeStreetMatch   bool
	PostalCodeStreetMatches []PostalCodeStreet
	StreetCityMatch         bool
	StreetCityMatches       []CityStreetPostalCode
	CityPostalCodeMatch     bool
	CityPostalCodeMatches   []CityPostalCode
	NominatimErrors         error
}

type CityPostalCode struct {
	City                string  `json:"city"`
	PostalCode          string  `json:"postal_code"`
	CountryCode         string  `json:"country_code"`
	CitySimilarity      float32 `json:"city_similarity"`
	WasPartialCityMatch bool    `json:"was_partial_city_match"`
	WasListMatch        bool    `json:"was_list_match"`
}

type CityStreetPostalCode struct {
	City                  string  `json:"city"`
	Street                string  `json:"street"`
	PostalCode            string  `json:"postal_code"`
	CountryCode           string  `json:"country_code"`
	StreetSimilarity      float32 `json:"street_similarity"`
	WasPartialStreetMatch bool    `json:"was_partial_street_match"`
	CitySimilarity        float32 `json:"city_similarity"`
	WasPartialCityMatch   bool    `json:"was_partial_city_match"`
	WasListMatch          bool    `json:"was_list_match"`
}

type PostalCodeStreet struct {
	PostalCode            string  `json:"postal_code"`
	Street                string  `json:"street"`
	CountryCode           string  `json:"country_code"`
	StreetSimilarity      float32 `json:"street_similarity"`
	WasPartialStreetMatch bool    `json:"was_partial_street_match"`
	WasListMatch          bool    `json:"was_list_match"`
}

type NominatimConfig struct {
	Url       string
	Languages []string
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CountryCodeCheck(countryCode string, matching PairMatching) bool {
	countryCode = strings.ToLower(countryCode)
	for _, val := range matching.CityPostalCodeMatches {
		if val.CountryCode == countryCode {
			return true
		}
	}
	for _, val := range matching.PostalCodeStreetMatches {
		if val.CountryCode == countryCode {
			return true
		}
	}
	for _, val := range matching.StreetCityMatches {
		if val.CountryCode == countryCode {
			return true
		}
	}
	return false
}
