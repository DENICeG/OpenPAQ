package nominatim_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"openPAQ/internal/algorithms"
	"openPAQ/internal/nominatim"
	"openPAQ/internal/normalization"
	"openPAQ/internal/types"
	"slices"
	"testing"

	"github.com/hbollon/go-edlib"
)

type ApiMock struct {
	searchStringResponse         []nominatim.NominatimCoreResult
	searchStringRaisesError      bool
	parameterResponse            []nominatim.NominatimCoreResult
	parameterResponseRaisesError bool
}

func (a ApiMock) ExecuteNominatimRequest(r *http.Request) ([]nominatim.NominatimCoreResult, error) {
	//TODO implement me
	panic("implement me")
}

func (a ApiMock) RequestWithSearchString(ctx context.Context, url, searchString, limit, language string) ([]nominatim.NominatimCoreResult, error) {

	if a.searchStringRaisesError {
		return []nominatim.NominatimCoreResult{}, errors.New("error")
	}
	return a.searchStringResponse, nil
}

func (a ApiMock) RequestWithParameters(ctx context.Context, url string, parameters nominatim.NominatimDetailRequest, limit string, language string) ([]nominatim.NominatimCoreResult, error) {

	if a.parameterResponseRaisesError {
		return []nominatim.NominatimCoreResult{}, errors.New("error")
	}

	return a.parameterResponse, nil
}

var singleResult = []nominatim.NominatimCoreResult{
	{
		City:        "CityOne",
		Road:        "StreetOne PartTwo",
		PostCode:    "12345",
		CountryCode: "DE",
	},
}

var multiResultAlmostSimilarCity = []nominatim.NominatimCoreResult{
	{
		City:        "CityOne",
		Road:        "StreetOne PartTwo",
		PostCode:    "12345",
		CountryCode: "DE",
	},
	{
		City:        "CityOne AdditionalName",
		Road:        "StreetOne PartTwo",
		PostCode:    "12345",
		CountryCode: "DE",
	},
}

var multiResultAlmostSimilarStreet = []nominatim.NominatimCoreResult{
	{
		City:        "CityOne",
		Road:        "StreetOne PartTwo",
		PostCode:    "12345",
		CountryCode: "DE",
	},
	{
		City:        "CityOne",
		Road:        "StreetOne PartTwoAndHalf",
		PostCode:    "12345",
		CountryCode: "DE",
	},
}

func comparePairmatching(first, second types.PairMatching) error {

	if first.CityPostalCodeMatch != second.CityPostalCodeMatch {
		return fmt.Errorf("CityStreetCheck().CityPostalCodeMatch = %v, want %v", first.CityPostalCodeMatch, second.CityPostalCodeMatch)
	}

	if len(first.CityPostalCodeMatches) != len(second.CityPostalCodeMatches) {
		return fmt.Errorf("CityStreetCheck().CityPostalCodeMatches = %v, want %v", first.CityPostalCodeMatches, second.CityPostalCodeMatches)
	}

	for i := range second.CityPostalCodeMatches {
		if !slices.Contains(first.CityPostalCodeMatches, second.CityPostalCodeMatches[i]) {
			return fmt.Errorf("CityStreetCheck().CityPostalCodeMatches = %v, want %v", first.CityPostalCodeMatches, second.CityPostalCodeMatches)
		}
	}

	if first.StreetCityMatch != second.StreetCityMatch {
		return fmt.Errorf("CityStreetCheck().StreetCityMatch = %v, want %v", first.StreetCityMatch, second.StreetCityMatch)
	}

	if len(first.StreetCityMatches) != len(second.StreetCityMatches) {
		return fmt.Errorf("CityStreetCheck().StreetCityMatches = %v, want %v", first.StreetCityMatches, second.StreetCityMatches)
	}

	for i := range second.StreetCityMatches {
		if !slices.Contains(first.StreetCityMatches, second.StreetCityMatches[i]) {
			return fmt.Errorf("CityStreetCheck().StreetCityMatches = %v, want %v", first.StreetCityMatches, second.StreetCityMatches)
		}
	}

	if first.PostalCodeStreetMatch != second.PostalCodeStreetMatch {
		return fmt.Errorf("CityStreetCheck().PostalCodeStreetMatch = %v, want %v", first.PostalCodeStreetMatch, second.PostalCodeStreetMatch)
	}

	if len(first.PostalCodeStreetMatches) != len(second.PostalCodeStreetMatches) {
		return fmt.Errorf("CityStreetCheck().StreetCityMatches = %v, want %v", first.PostalCodeStreetMatches, second.PostalCodeStreetMatches)
	}

	for i := range second.PostalCodeStreetMatches {
		if !slices.Contains(first.PostalCodeStreetMatches, second.PostalCodeStreetMatches[i]) {
			return fmt.Errorf("CityStreetCheck().StreetCityMatches = %v, want %v", first.PostalCodeStreetMatches, second.PostalCodeStreetMatches)
		}
	}

	return nil
}

func TestCityStreetCheck(t *testing.T) {
	config := algorithms.MatchSeverityConfig{
		Algorithm:          edlib.Lcs,
		AlgorithmThreshold: 0.65,
	}

	tests := []struct {
		nom   nominatim.Nominatim
		name  string
		input types.NormalizeInput
		want  types.PairMatching
	}{
		{
			name: "One response from nominatim - happy case",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringResponse: singleResult,
					parameterResponse:    singleResult,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			want: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "cityone",
						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
						CitySimilarity:        1,
						WasPartialCityMatch:   false,
					},
				},
			},
		},
		{
			name: "multi response from nominatim - happy case",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringResponse: singleResult,
					parameterResponse:    multiResultAlmostSimilarCity,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			want: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "cityone",
						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
						CitySimilarity:        1,
						WasPartialCityMatch:   false,
					}, {
						City:                  "cityone additionalname",
						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
						CitySimilarity:        0.3181818,
						WasPartialCityMatch:   true,
					},
				},
			},
		},
		{
			name: "RequestWithSearchString - raise error",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringRaisesError: true,
					parameterResponse:       singleResult,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			want: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "cityone",
						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
						CitySimilarity:        1,
						WasPartialCityMatch:   false,
					},
				},
			},
		},
		{
			name: "RequestWithParameter - raise error",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					parameterResponseRaisesError: true,
					searchStringResponse:         singleResult,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			want: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "cityone",
						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
						CitySimilarity:        1,
						WasPartialCityMatch:   false,
					},
				},
			},
		},
		{
			name: "RequestWithParameter and SearchString raise error",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config,
				normalization.NewNormalizer("generic"),
				ApiMock{
					parameterResponseRaisesError: true,
					searchStringRaisesError:      true,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			want: types.PairMatching{},
		},
		{
			name: "Partial Match - 1",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringResponse: singleResult,
					parameterResponse:    singleResult,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttw"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "pl",
			},
			want: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "cityone",
						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      0.9411765,
						WasPartialStreetMatch: false,
						CitySimilarity:        1,
						WasPartialCityMatch:   false,
					},
				},
			},
		},
		{
			name: "Partial Match - 2",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringResponse: singleResult,
					parameterResponse:    singleResult,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "pl",
			},
			want: types.PairMatching{
				StreetCityMatch: true,
				StreetCityMatches: []types.CityStreetPostalCode{
					{
						City:                  "cityone",
						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      0.5294118,
						WasPartialStreetMatch: true,
						CitySimilarity:        1,
						WasPartialCityMatch:   false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := comparePairmatching(<-tt.nom.CityStreetCheck(context.Background(), tt.input), tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestPostalCodeStreetCheck(t *testing.T) {
	config := algorithms.MatchSeverityConfig{
		Algorithm:          edlib.Lcs,
		AlgorithmThreshold: 0.65,
	}

	tests := []struct {
		name              string
		nom               nominatim.Nominatim
		input             types.NormalizeInput
		cityStreetMatches []types.CityStreetPostalCode
		want              types.PairMatching
	}{
		{
			name: "cityStreetMatches contains postal code and street",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			cityStreetMatches: []types.CityStreetPostalCode{
				{
					City:                  "cityone",
					Street:                "streetone parttwo",
					PostalCode:            "12345",
					CountryCode:           "de",
					StreetSimilarity:      1,
					WasPartialStreetMatch: false,
					CitySimilarity:        1,
					WasPartialCityMatch:   false,
				},
			},
			want: types.PairMatching{
				PostalCodeStreetMatch: true,
				PostalCodeStreetMatches: []types.PostalCodeStreet{
					{

						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
					},
				},
			},
		},
		{
			name: "happy case",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringResponse: singleResult,
					parameterResponse:    singleResult,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			cityStreetMatches: nil,
			want: types.PairMatching{
				PostalCodeStreetMatch: true,
				PostalCodeStreetMatches: []types.PostalCodeStreet{
					{

						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      1,
						WasPartialStreetMatch: false,
					},
				},
			},
		},
		{
			name: "partial street poland",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringResponse: singleResult,
					parameterResponse:    multiResultAlmostSimilarStreet,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "pl",
			},
			cityStreetMatches: nil,
			want: types.PairMatching{
				PostalCodeStreetMatch: true,
				PostalCodeStreetMatches: []types.PostalCodeStreet{
					{

						Street:                "streetone parttwo",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      0.5294118,
						WasPartialStreetMatch: true,
					},
					{

						Street:                "streetone parttwoandhalf",
						PostalCode:            "12345",
						CountryCode:           "de",
						StreetSimilarity:      0.375,
						WasPartialStreetMatch: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := comparePairmatching(<-tt.nom.PostalCodeStreetCheck(context.Background(), tt.input, tt.cityStreetMatches), tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestPostalCodeCityCheck(t *testing.T) {
	config := algorithms.MatchSeverityConfig{
		Algorithm:          edlib.Lcs,
		AlgorithmThreshold: 0.65,
	}

	tests := []struct {
		name              string
		nom               nominatim.Nominatim
		input             types.NormalizeInput
		cityStreetMatches []types.CityStreetPostalCode
		want              types.PairMatching
	}{
		{
			name: "cityStreetMatches contains postal code and street",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			cityStreetMatches: []types.CityStreetPostalCode{
				{
					City:                  "cityone",
					Street:                "streetone parttwo",
					PostalCode:            "12345",
					CountryCode:           "de",
					StreetSimilarity:      1,
					WasPartialStreetMatch: false,
					CitySimilarity:        1,
					WasPartialCityMatch:   false,
				},
			},
			want: types.PairMatching{
				CityPostalCodeMatch: true,
				CityPostalCodeMatches: []types.CityPostalCode{
					{
						City:                "cityone",
						PostalCode:          "12345",
						CountryCode:         "de",
						CitySimilarity:      1,
						WasPartialCityMatch: false,
					},
				},
			},
		},
		{
			name: "happy case",
			nom: *nominatim.NewNominatim(
				"",
				[]string{"en", "de"},
				config, normalization.NewNormalizer("generic"),
				ApiMock{
					searchStringResponse: singleResult,
					parameterResponse:    singleResult,
				},
				"test",
			),
			input: types.NormalizeInput{
				Streets:     []string{"streetone parttwo"},
				City:        "cityone",
				PostalCode:  "12345",
				CountryCode: "de",
			},
			cityStreetMatches: nil,
			want: types.PairMatching{
				CityPostalCodeMatch: true,
				CityPostalCodeMatches: []types.CityPostalCode{
					{
						City:                "cityone",
						PostalCode:          "12345",
						CountryCode:         "de",
						CitySimilarity:      1,
						WasPartialCityMatch: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := comparePairmatching(<-tt.nom.PostalCodeCityCheck(context.Background(), tt.input, tt.cityStreetMatches), tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}
