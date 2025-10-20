package main

import (
	"fmt"
	"openPAQ/internal"
	"openPAQ/internal/algorithms"
	diytypes "openPAQ/internal/listmatcher/types"
	"openPAQ/internal/types"
	"os"
	"strings"

	"github.com/hbollon/go-edlib"
	log "github.com/sirupsen/logrus"
)

func init() {
	loglevel := lookupEnv("LOG_LEVEL")
	level, err := log.ParseLevel(loglevel)
	if err != nil {
		setupLogger(log.DebugLevel)
	} else {
		setupLogger(level)
	}
}

func setupLogger(logLevel log.Level) {
	// Log as JSON instead of the default ASCII formatter.

	if os.Getenv("LOG_FORMAT") == "text" {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(logLevel)
}

func lookupEnv(s string) string {
	if e, ok := os.LookupEnv(s); !ok {
		panic(fmt.Errorf("ENVVAR: %s is not set", s))
	} else {
		return e
	}
}

func main() {

	clickhouseEnabledString := lookupEnv("CLICKHOUSE_ENABLED")
	clickhouseEnabled := false
	if clickhouseEnabledString == "true" || clickhouseEnabledString == "1" {
		clickhouseEnabled = true
	}

	var databaseConfig diytypes.DatabaseConfig
	if clickhouseEnabled == true {
		databaseConfig = diytypes.DatabaseConfig{
			DbUserName:     lookupEnv("CLICKHOUSE_DB_USERNAME"),
			DbUserPassword: lookupEnv("CLICKHOUSE_DB_PASSWORD"),
			DbHost:         lookupEnv("CLICKHOUSE_DB_HOST"),
			DbPort:         lookupEnv("CLICKHOUSE_DB_PORT"),
			DataBase:       lookupEnv("CLICKHOUSE_DB_DATABASE"),
			Table:          lookupEnv("CLICKHOUSE_DB_TABLE"),
		}
	}
	matchSeverityConfig := algorithms.MatchSeverityConfig{
		Algorithm:                 edlib.DamerauLevenshtein,
		AlgorithmThreshold:        0.8,
		PartialAlgorithm:          edlib.DamerauLevenshtein,
		PartialAlgorithmThreshold: 0.8,
	}

	useTlsEnv := lookupEnv("USE_TLS")
	var useTls bool
	var tlsKeyFilePath string
	var tlsCertFilePath string

	if useTlsEnv == "1" || useTlsEnv == "true" {
		useTls = true
		tlsKeyFilePath = lookupEnv("TLS_KEY_FILE_PATH")
		tlsCertFilePath = lookupEnv("TLS_CERT_FILE_PATH")
	}

	useJwtEnv := lookupEnv("USE_JWT")
	var useJwt bool
	var jwtSigningKey []byte
	if useJwtEnv == "1" || useJwtEnv == "true" {
		useJwt = true
		jwtSigningKey = []byte(lookupEnv("JWT_SIGNING_KEY"))

	}

	webserverConfig := internal.WebserverConfig{
		ListenAddress:   lookupEnv("WEBSERVER_LISTEN_ADDRESS"),
		JwtSigningKey:   jwtSigningKey,
		UseJWT:          useJwt,
		UseTLS:          useTls,
		TLSKeyFilePath:  tlsKeyFilePath,
		TLSCertFilePath: tlsCertFilePath,
	}

	nominatimConfig := types.NominatimConfig{
		Url:       lookupEnv("NOMINATIM_ADDRESS"),
		Languages: []string{"de", "en"},
	}

	ce := lookupEnv("CACHE_ENABLED")
	enableCache := false
	cacheUrl := ""
	if strings.ToLower(ce) == "true" || ce == "1" {
		enableCache = true
		cacheUrl = lookupEnv("CACHE_URL")
	}

	service := internal.NewService(&internal.ServiceConfig{
		Webserver:         webserverConfig,
		DIYDatabaseConfig: databaseConfig,
		Version:           lookupEnv("VERSION"),
		UseCaching:        enableCache,
		CacheUrl:          cacheUrl,
		ClickhouseEnabled: clickhouseEnabled,
	}, matchSeverityConfig, nominatimConfig)

	log.Error(service.Start())
}
