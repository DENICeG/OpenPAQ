package internal

import (
	"fmt"
	"net/http"
	"openPAQ/internal/algorithms"
	"openPAQ/internal/listmatcher"
	types2 "openPAQ/internal/listmatcher/types"
	"openPAQ/internal/nominatim"
	"openPAQ/internal/normalization"
	"openPAQ/internal/types"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

type ServiceConfig struct {
	Webserver         WebserverConfig
	DIYDatabaseConfig types2.DatabaseConfig
	Version           string
	CacheUrl          string
	UseCaching        bool
	ClickhouseEnabled bool
}

type Service struct {
	engine      *gin.Engine
	webserver   *http.Server
	config      *ServiceConfig
	listMatcher *listmatcher.ListMatcher
	nominatim   *nominatim.Nominatim
	normalizer  *normalization.Normalizer
	mc          *memcache.Client
}

func NewService(config *ServiceConfig, matcherConfig algorithms.MatchSeverityConfig, nominatimConfig types.NominatimConfig) *Service {
	var d *listmatcher.ListMatcher
	if config.ClickhouseEnabled {
		d = listmatcher.NewMatcher(matcherConfig)

		if err := d.Register("de", config.DIYDatabaseConfig, matcherConfig); err != nil {
			panic("unable to register DE country checker")
		}
	}

	var mc *memcache.Client
	if config.UseCaching {
		mc = memcache.New(config.CacheUrl)
	}

	norma := normalization.NewNormalizer("generic")

	service := &Service{
		engine:      gin.New(),
		webserver:   nil,
		config:      config,
		listMatcher: d,
		nominatim:   nominatim.NewNominatim(nominatimConfig.Url, nominatimConfig.Languages, matcherConfig, norma, nil, fmt.Sprint("openPAQ", "-", config.Version)),
		normalizer:  norma,
		mc:          mc,
	}

	service.webserver = &http.Server{
		Addr:         config.Webserver.ListenAddress,
		Handler:      service.engine,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	service.setupWebserver()
	return service
}

func (s *Service) Start() error {
	if s.config.Webserver.UseTLS {
		return s.webserver.ListenAndServeTLS(s.config.Webserver.TLSCertFilePath, s.config.Webserver.TLSKeyFilePath)
	}
	return s.webserver.ListenAndServe()

}
