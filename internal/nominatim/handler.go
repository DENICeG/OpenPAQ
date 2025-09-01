package nominatim

import (
	"context"
	"net/http"
	"openPAQ/internal/algorithms"
	"openPAQ/internal/normalization"
	"openPAQ/internal/types"
	"time"
)

type Nominatim struct {
	url        string
	languages  []string
	config     algorithms.MatchSeverityConfig
	api        api
	normalizer *normalization.Normalizer
}

func NewNominatim(url string, languages []string, config algorithms.MatchSeverityConfig, normalizer *normalization.Normalizer, nominatimApi api, userAgent string) *Nominatim {

	nominatim := Nominatim{
		url,
		languages,
		config,
		apiNominatim{
			client: http.Client{
				Timeout: 180 * time.Second,
			},
			userAgent: userAgent,
		},
		normalizer,
	}

	if nominatimApi != nil {
		nominatim.api = nominatimApi
	}
	return &nominatim
}

func (nom *Nominatim) Handle(ctx context.Context, input types.Input) <-chan types.PairMatching {

	c := make(chan types.PairMatching, 1)

	go func() {
		defer close(c)
		var result types.PairMatching

		normalizedInput := input.Normalize()

		rC := nom.CityStreetCheck(ctx, normalizedInput)
		r := <-rC

		result.StreetCityMatch = r.StreetCityMatch
		result.StreetCityMatches = r.StreetCityMatches

		pC := nom.PostalCodeCityCheck(ctx, normalizedInput, r.StreetCityMatches)
		qC := nom.PostalCodeStreetCheck(ctx, normalizedInput, r.StreetCityMatches)

		for counter := 0; counter < 2; counter++ {
			select {
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
