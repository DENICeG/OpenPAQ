package internal

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"openPAQ/internal/algorithms"
	"openPAQ/internal/types"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	activeRequest = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "datascience",
		Name:      "activeRequests",
	}, []string{"user"})
	allRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "datascience",
		Name:      "processed_requests",
	}, []string{"user"})
)

type WebserverConfig struct {
	JwtSigningKey   []byte
	ListenAddress   string
	TLSKeyFilePath  string
	TLSCertFilePath string
	UseTLS          bool
	UseJWT          bool
}

func (s *Service) setupWebserver() {
	s.engine.Use(gin.Recovery())
	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	s.engine.GET("/version", func(context *gin.Context) {
		context.JSON(200, gin.H{"version": s.config.Version})
	})

	if s.config.Webserver.UseJWT {
		s.engine.GET("/api/v1/check", s.CheckJwt(), s.checkHandler)
	} else {
		s.engine.GET("/api/v1/check", s.checkHandler)
	}

	prometheus.MustRegister(types.InputNormalizerErrorCounter, activeRequest, allRequests)

	s.engine.GET("/metrics", func(ctx *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(ctx.Writer, ctx.Request)
	})

	s.engine.LoadHTMLFiles("./docs/docs/openAPI/index.html")
	s.engine.Static("/api/doc", "./docs/docs/openAPI")

	if err := s.startWebserver(); err != nil {
		log.Warn(err)
	}
}

func (s *Service) startWebserver() error {
	if s.config.Webserver.UseTLS {
		err := s.webserver.ListenAndServeTLS(s.config.Webserver.TLSCertFilePath, s.config.Webserver.TLSKeyFilePath)
		if err != nil {
			return err
		}
	} else {
		err := s.webserver.ListenAndServe()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) checkHandler(ctx *gin.Context) {
	startTime := time.Now()
	debugDetails := ctx.Query("debug_details")
	street := ctx.Query("street")
	city := ctx.Query("city")
	postalCode := ctx.Query("postal_code")
	countryCode := ctx.Query("country_code")

	if debugDetails != "true" {
		debugDetails = "false"
	}

	if len(street) > 500 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "street field exceed length limit of 500 elements",
		})
		return
	}

	if len(city) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "city field exceed length limit of 100 elements",
		})
		return
	}

	if len(postalCode) > 50 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "postal code field exceed length limit of 50 elements",
		})
		return
	}

	if len(countryCode) > 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "country code field exceed length limit of 2 elements",
		})
		return
	}

	input := types.Input{
		Street:      street,
		City:        city,
		PostalCode:  postalCode,
		CountryCode: strings.ToLower(countryCode),
		Normalizer:  s.normalizer.Get(strings.ToLower(countryCode)),
	}

	defer func() {
		if val, ok := ctx.Get("SubjectFromToken"); ok {
			activeRequest.WithLabelValues(val.(string)).Dec()
			allRequests.WithLabelValues(val.(string)).Inc()
		}
	}()

	if val, ok := ctx.Get("SubjectFromToken"); ok {
		activeRequest.WithLabelValues(val.(string)).Inc()
	}

	if s.config.UseCaching {
		cacheResult, err := s.checkCache(input)

		var res types.Result
		res = cacheResult.Result
		res.Street = street
		res.City = city
		res.PostalCode = postalCode
		res.CountryCode = countryCode

		if err == nil {
			if debugDetails == "true" {
				var matcherConfig algorithms.MatchSeverityConfig
				if s.config.ClickhouseEnabled {
					matcherConfig = s.listMatcher.GetConfig()
				}
				debugRes := types.DebugResult{
					Result: res,
					DebugDetails: types.DebugDetails{
						MatchSeverityConfig:     matcherConfig,
						StreetCityMatches:       cacheResult.StreetCityMatches,
						StreetPostalCodeMatches: cacheResult.StreetPostalCodeMatches,
						PostalCodeCityMatches:   cacheResult.PostalCodeCityMatches,
					},
				}
				ctx.JSON(http.StatusOK, debugRes)
				log.Debugf("got info from cache")
				return
			} else {
				ctx.JSON(http.StatusOK, res)
				log.Debugf("got info from cache")
				return
			}
		}

		if !errors.Is(err, memcache.ErrCacheMiss) {
			log.WithFields(log.Fields{"error": err}).Error("error communicating with cache")
		}
	}

	var source types.SourceOfTruth
	var debugDude types.PairMatching
	var NominatimError error

	ctx2, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if s.config.ClickhouseEnabled && s.listMatcher.Possible(input.CountryCode) {
		func() {
			defer cancel()

			pairMatchesDiy := s.listMatcher.Handle(ctx2, input)
			pairMatchesNominatim := s.nominatim.Handle(ctx2, input)

			for counter := 0; counter < 2; counter++ {
				select {
				case r := <-pairMatchesDiy:
					pairMatchesDiy = nil
					if res, ok := eval(r, input.Normalize().CountryCode); ok {
						debugDude = r
						source = res
						return
					} else {
						debugDude.CityPostalCodeMatches = append(debugDude.CityPostalCodeMatches, r.CityPostalCodeMatches...)
						debugDude.PostalCodeStreetMatches = append(debugDude.PostalCodeStreetMatches, r.PostalCodeStreetMatches...)
						debugDude.StreetCityMatches = append(debugDude.StreetCityMatches, r.StreetCityMatches...)
						source.PostalCodeMatched = source.PostalCodeMatched || res.PostalCodeMatched
						source.CityToPostalCodeMatched = source.CityToPostalCodeMatched || res.CityToPostalCodeMatched
						source.CountryCodeMatched = source.CountryCodeMatched || res.CountryCodeMatched
						source.CityMatched = source.CityMatched || res.CityMatched
						source.StreetMatched = source.StreetMatched || res.StreetMatched
					}
				case r := <-pairMatchesNominatim:
					pairMatchesNominatim = nil
					if r.NominatimErrors != nil {
						NominatimError = r.NominatimErrors
						return
					}
					if res, ok := eval(r, input.CountryCode); ok {
						debugDude = r
						source = res
						return
					} else {
						debugDude.CityPostalCodeMatches = append(debugDude.CityPostalCodeMatches, r.CityPostalCodeMatches...)
						debugDude.PostalCodeStreetMatches = append(debugDude.PostalCodeStreetMatches, r.PostalCodeStreetMatches...)
						debugDude.StreetCityMatches = append(debugDude.StreetCityMatches, r.StreetCityMatches...)
						source.PostalCodeMatched = source.PostalCodeMatched || res.PostalCodeMatched
						source.CityToPostalCodeMatched = source.CityToPostalCodeMatched || res.CityToPostalCodeMatched
						source.CountryCodeMatched = source.CountryCodeMatched || res.CountryCodeMatched
						source.CityMatched = source.CityMatched || res.CityMatched
						source.StreetMatched = source.StreetMatched || res.StreetMatched
					}
				case <-ctx2.Done():
					return
				}
			}
		}()
	} else {
		pairMatchesNominatim := <-s.nominatim.Handle(ctx2, input)
		if pairMatchesNominatim.NominatimErrors != nil {
			NominatimError = pairMatchesNominatim.NominatimErrors
		}
		source.StreetMatched = pairMatchesNominatim.StreetCityMatch || pairMatchesNominatim.PostalCodeStreetMatch
		source.CityMatched = pairMatchesNominatim.CityPostalCodeMatch || pairMatchesNominatim.StreetCityMatch
		source.PostalCodeMatched = pairMatchesNominatim.PostalCodeStreetMatch || pairMatchesNominatim.CityPostalCodeMatch
		source.CityToPostalCodeMatched = pairMatchesNominatim.CityPostalCodeMatch
		source.CountryCodeMatched = types.CountryCodeCheck(input.CountryCode, pairMatchesNominatim)
		debugDude = pairMatchesNominatim
	}

	if NominatimError != nil {

		msg := strings.Split(NominatimError.Error(), "\n")

		log.WithFields(log.Fields{
			"msg": msg,
		}).Errorf("Request to Nominatim failed")

		ctx.JSON(500, gin.H{
			"NominatimError": msg,
		})
		return
	}

	var matcherConfig algorithms.MatchSeverityConfig
	if s.config.ClickhouseEnabled {
		matcherConfig = s.listMatcher.GetConfig()
	}
	result := types.DebugResult{
		Result: types.Result{
			Input:         input,
			SourceOfTruth: source,
			Version:       s.config.Version,
		},
		DebugDetails: types.DebugDetails{
			MatchSeverityConfig:     matcherConfig,
			StreetCityMatches:       types.RemoveDuplicate(append(debugDude.StreetCityMatches)),
			StreetPostalCodeMatches: types.RemoveDuplicate(append(debugDude.PostalCodeStreetMatches)),
			PostalCodeCityMatches:   types.RemoveDuplicate(append(debugDude.CityPostalCodeMatches)),
		},
	}
	if s.config.UseCaching {
		if err := s.putCache(input, result); err != nil {
			log.WithFields(log.Fields{"error": err}).Errorf("error storing result in cache")
		}
	}

	log.WithFields(log.Fields{"duration": time.Since(startTime).Milliseconds()}).Debugf("duration needed for requestWithSearchString")

	if debugDetails != "true" {
		ctx.JSON(http.StatusOK, result.Result)
	} else {
		ctx.JSON(http.StatusOK, result)
	}

}

func (s *Service) CheckJwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.Request.Header.Get("Authorization")
		token = strings.ReplaceAll(token, "Bearer ", "")

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.config.Webserver.JwtSigningKey, nil
		})

		if parsedToken == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "No token provided",
			})
			ctx.Abort()
			return
		}

		switch {
		case parsedToken.Valid:
			subject, err := parsedToken.Claims.GetSubject()

			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Invalid token",
				})
				ctx.Abort()
			}
			ctx.Set("SubjectFromToken", subject)
			ctx.Next()
			break
		case errors.Is(err, jwt.ErrTokenMalformed):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid token",
			})
			ctx.Abort()

		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token signature",
			})
			ctx.Abort()

		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token expired",
			})

			sub, err := parsedToken.Claims.GetSubject()

			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Invalid token",
				})
				ctx.Abort()
			}
			log.WithFields(log.Fields{
				"Subject": sub,
			}).Error("Token expired")

			ctx.Abort()
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not parse token",
			})
			ctx.Abort()
		}
	}
}

func eval(input types.PairMatching, cc string) (result types.SourceOfTruth, ok bool) {
	result.StreetMatched = input.StreetCityMatch || input.PostalCodeStreetMatch
	result.CityMatched = input.CityPostalCodeMatch || input.StreetCityMatch
	result.PostalCodeMatched = input.PostalCodeStreetMatch || input.CityPostalCodeMatch
	result.CityToPostalCodeMatched = input.CityPostalCodeMatch

	for _, v := range input.CityPostalCodeMatches {
		if cc == v.CountryCode {
			result.CountryCodeMatched = true
			break
		}
	}

	if !result.CountryCodeMatched {
		for _, v := range input.StreetCityMatches {
			if cc == v.CountryCode {
				result.CountryCodeMatched = true
				break
			}
		}
	}

	if !result.CountryCodeMatched {
		for _, v := range input.PostalCodeStreetMatches {
			if cc == v.CountryCode {
				result.CountryCodeMatched = true
				break
			}
		}
	}

	if result.StreetMatched && result.CityMatched && result.PostalCodeMatched && result.CityToPostalCodeMatched && result.CountryCodeMatched {
		ok = true
	} else {
		ok = false
	}

	return
}

func (s *Service) checkCache(input types.Input) (types.DebugResult, error) {

	res, err := s.mc.Get(s.createHash(input))
	if err != nil {
		return types.DebugResult{}, err
	}

	var result types.DebugResult
	if err = json.Unmarshal(res.Value, &result); err != nil {
		return types.DebugResult{}, err
	}

	return result, nil
}

func (s *Service) putCache(input types.Input, res types.DebugResult) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return s.mc.Add(&memcache.Item{
		Key:   s.createHash(input),
		Value: data,
	})
}

func (s *Service) createHash(input types.Input) string {
	i := input.Normalize()
	m := sha1.Sum([]byte(fmt.Sprintf("%s%s%s%s%s", i.Streets, i.City, i.PostalCode, i.CountryCode, s.config.Version)))
	return hex.EncodeToString(m[:])
}
