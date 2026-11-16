// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"openPAQ/internal/types"
	"reflect"
	"testing"
	"time"
)

func Test_eval(t *testing.T) {
	type args struct {
		input types.PairMatching
		cc    string
	}
	tests := []struct {
		name       string
		args       args
		wantResult types.SourceOfTruth
		wantOk     bool
	}{
		{
			name: "perfect score",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   true,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         true,
					StreetCityMatches:       nil,
					CityPostalCodeMatch:     true,
					CityPostalCodeMatches: []types.CityPostalCode{{
						City:        "",
						PostalCode:  "",
						CountryCode: "de",
					}},
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             true,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: true,
				CountryCodeMatched:      true,
			},
			wantOk: true,
		},
		{
			name: "very bad",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   false,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         false,
					StreetCityMatches:       nil,
					CityPostalCodeMatch:     false,
					CityPostalCodeMatches: []types.CityPostalCode{{
						City:        "",
						PostalCode:  "",
						CountryCode: "dk",
					}},
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           false,
				CityMatched:             false,
				PostalCodeMatched:       false,
				CityToPostalCodeMatched: false,
				CountryCodeMatched:      false,
			},
			wantOk: false,
		},
		{
			name: "wrong country",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   true,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         true,
					StreetCityMatches:       nil,
					CityPostalCodeMatch:     true,
					CityPostalCodeMatches: []types.CityPostalCode{{
						City:        "",
						PostalCode:  "",
						CountryCode: "dk",
					}},
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             true,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: true,
				CountryCodeMatched:      false,
			},
			wantOk: false,
		},
		{
			name: "wrong country, street",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   false,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         false,
					StreetCityMatches:       nil,
					CityPostalCodeMatch:     true,
					CityPostalCodeMatches: []types.CityPostalCode{{
						City:        "",
						PostalCode:  "",
						CountryCode: "dk",
					}},
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           false,
				CityMatched:             true,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: true,
				CountryCodeMatched:      false,
			},
			wantOk: false,
		},
		{
			name: "correct country, wrong street",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   false,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         false,
					StreetCityMatches:       nil,
					CityPostalCodeMatch:     true,
					CityPostalCodeMatches: []types.CityPostalCode{{
						City:        "",
						PostalCode:  "",
						CountryCode: "de",
					}},
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           false,
				CityMatched:             true,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: true,
				CountryCodeMatched:      true,
			},
			wantOk: false,
		},
		{
			name: "wrong country, city  <-> postalcode",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch: true,
					PostalCodeStreetMatches: []types.PostalCodeStreet{{
						PostalCode:  "",
						Street:      "",
						CountryCode: "dk",
					}},
					StreetCityMatch:       true,
					StreetCityMatches:     nil,
					CityPostalCodeMatch:   false,
					CityPostalCodeMatches: nil,
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             true,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: false,
				CountryCodeMatched:      false,
			},
			wantOk: false,
		},
		{
			name: "correct country, wrong city <-> postalcode",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   true,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         true,
					StreetCityMatches: []types.CityStreetPostalCode{{
						City:        "",
						Street:      "",
						PostalCode:  "",
						CountryCode: "de",
					}},
					CityPostalCodeMatch:   false,
					CityPostalCodeMatches: nil,
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             true,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: false,
				CountryCodeMatched:      true,
			},
			wantOk: false,
		},
		{
			name: "wrong country, wrong postalcode",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   false,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         true,
					StreetCityMatches: []types.CityStreetPostalCode{{
						City:        "",
						Street:      "",
						PostalCode:  "",
						CountryCode: "dk",
					}},
					CityPostalCodeMatch:   false,
					CityPostalCodeMatches: nil,
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             true,
				PostalCodeMatched:       false,
				CityToPostalCodeMatched: false,
				CountryCodeMatched:      false,
			},
			wantOk: false,
		},
		{
			name: "correct country, wrong postalcode",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch:   false,
					PostalCodeStreetMatches: nil,
					StreetCityMatch:         true,
					StreetCityMatches: []types.CityStreetPostalCode{{
						City:        "",
						Street:      "",
						PostalCode:  "",
						CountryCode: "de",
					}},
					CityPostalCodeMatch:   false,
					CityPostalCodeMatches: nil,
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             true,
				PostalCodeMatched:       false,
				CityToPostalCodeMatched: false,
				CountryCodeMatched:      true,
			},
			wantOk: false,
		},
		{
			name: "wrong country, wrong city",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch: true,
					PostalCodeStreetMatches: []types.PostalCodeStreet{{
						PostalCode:  "",
						Street:      "",
						CountryCode: "dk",
					}},
					StreetCityMatch:       false,
					StreetCityMatches:     nil,
					CityPostalCodeMatch:   false,
					CityPostalCodeMatches: nil,
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             false,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: false,
				CountryCodeMatched:      false,
			},
			wantOk: false,
		},
		{
			name: "correct country, wrong city",
			args: args{
				input: types.PairMatching{
					PostalCodeStreetMatch: true,
					PostalCodeStreetMatches: []types.PostalCodeStreet{{
						PostalCode:  "",
						Street:      "",
						CountryCode: "de",
					}},
					StreetCityMatch:       false,
					StreetCityMatches:     nil,
					CityPostalCodeMatch:   false,
					CityPostalCodeMatches: nil,
				},
				cc: "de",
			},
			wantResult: types.SourceOfTruth{
				StreetMatched:           true,
				CityMatched:             false,
				PostalCodeMatched:       true,
				CityToPostalCodeMatched: false,
				CountryCodeMatched:      true,
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := eval(tt.args.input, tt.args.cc)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("eval() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("eval() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestTLSOff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gengine := gin.Default()
	gengine.GET("/api/v1/check", func(context *gin.Context) {
		context.Writer.Write([]byte("greetings"))
	})
	s := Service{
		engine: gengine,
		webserver: &http.Server{
			Addr:    ":25333",
			Handler: gengine,
		},
		config: &ServiceConfig{
			Webserver: WebserverConfig{
				JwtSigningKey:   nil,
				ListenAddress:   "",
				TLSKeyFilePath:  "",
				TLSCertFilePath: "",
				UseTLS:          false,
				UseJWT:          false,
			},
		},
	}

	go func() {
		if err := s.startWebserver(); err != nil {
			t.Error(err)
		}
	}()

	for i := 0; i < 3; i++ {
		<-time.After(3 * time.Second)
		resp, err := http.Get("http://localhost:25333/api/v1/check")
		if err != nil {
			t.Error("http request failed")
		}
		if resp == nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("http request returned wrong status code, expected 200, got: %d", resp.StatusCode)
		}
		break
	}
}

func TestTLSOn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gengine := gin.Default()

	s := Service{
		engine: gengine,
		webserver: &http.Server{
			Addr:    ":25334",
			Handler: gengine,
		},
		config: &ServiceConfig{
			Webserver: WebserverConfig{
				JwtSigningKey:   []byte("extremelysecret"),
				ListenAddress:   "",
				TLSKeyFilePath:  "",
				TLSCertFilePath: "",
				UseTLS:          true,
				UseJWT:          false,
			},
		},
	}

	if err := s.startWebserver(); err == nil {
		t.Error("did not fail loading files")
	}
}
