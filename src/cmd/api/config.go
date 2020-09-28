package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func stringToBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "1", "true":
		return true, nil
	case "0", "false":
		return false, nil
	default:
		return false, fmt.Errorf("cast '%s' to bool failed", s)
	}
}

func stringToBoolFallback(s string, fallback bool) bool {
	value, err := stringToBool(s)
	if err != nil {
		return fallback
	}
	return value
}

var suffixes = map[string]time.Duration{
	"ms": time.Millisecond,
	"s":  time.Second,
	"m":  time.Minute,
	"h":  time.Hour,
	"d":  time.Second,
}

func stringToDuration(s string) (time.Duration, error) {
	fields := strings.Fields(s)

	var number int
	var unit time.Duration

	switch len(fields) {
	case 1: // value and unit without whitespace separation
		input := fields[0]
		var value, suffix string
		var valEnded int

		for i, r := range input {
			_, err := strconv.Atoi(string(r))
			if err != nil {
				valEnded = i
				break
			}
			value += string(r)
		}
		if valEnded == 0 {
			return 0, fmt.Errorf("no value number provided")
		}
		var err error
		number, err = strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("cast value failed: %s", err)
		}

		suffix = input[valEnded:]
		for s, d := range suffixes {
			if suffix == s {
				unit = d
				break
			}
		}
	case 2: // value and unit with whitespace separation
		value, suffix := fields[0], fields[1]

		var err error
		number, err = strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("can't convert %s to an integer", value)
		}

		for s, d := range suffixes {
			if s == suffix {
				unit = d
				break
			}
		}

	default:
		return 0, fmt.Errorf("unexpected duration format, expected 2 parts at maximum")
	}

	if unit == 0 {
		var supportedList []string
		for s := range suffixes {
			supportedList = append(supportedList, s)
		}
		supported := strings.Join(supportedList, ", ")
		return 0, fmt.Errorf("provided suffix not supported (%s)", supported)
	}
	return time.Duration(int(unit) * number), nil
}

func stringFallback(s string, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func stringToDurationFallback(s string, fallback time.Duration) time.Duration {
	value, err := stringToDuration(s)
	if err != nil {
		return fallback
	}
	return value
}

func stringToIntFallback(s string, fallback int) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return value
}

type Config struct {
	serverReadTimeout  time.Duration
	serverWriteTimeout time.Duration

	answerCache                bool // enable/disable cache based on query sent to application and it's answer
	answerCacheExpiration      time.Duration
	answerCacheCleanupInterval time.Duration

	listingCache                bool // enable/disable cache for listing metadata for given "title" part of query
	listingCacheExpiration      time.Duration
	listingCacheCleanupInterval time.Duration

	contentCache                bool // enable/disable cache for downloaded content (strongly suggested)
	contentCacheExpiration      time.Duration
	contentCacheCleanupInterval time.Duration

	downloadDelayMin time.Duration // download Task limitation to emulate human-like behaviour to prevent from banning.
	downloadDelayMax time.Duration // min/max value, each download gets random value from that range

	searchWorkers      int           // search worker goroutines (inefficient without cached content)
	searchMaxDistance  int           // fuzzy-search engine distance option
	searchRandomResult bool          // returns a random match in the scope of given book instead of a first found match [Note: cannot work properly with with CACHE_ANSWER enabled]
	searchTimeout      time.Duration // maximum time allowed to spent by server for each search request

	providerUserAgent string        // user-agent header used for provider's requests
	providerTimeout   time.Duration // provider http client timeout
}

func GetDefaultConfig() *Config {
	return &Config{
		serverReadTimeout:  time.Minute * 2,
		serverWriteTimeout: time.Minute * 2,

		answerCache:                true,
		answerCacheExpiration:      time.Hour * 4,
		answerCacheCleanupInterval: time.Minute * 31,

		listingCache:                true,
		listingCacheExpiration:      time.Hour * 4,
		listingCacheCleanupInterval: time.Minute * 10,

		contentCache:                true,
		contentCacheExpiration:      time.Hour,
		contentCacheCleanupInterval: time.Minute * 10,

		downloadDelayMin: time.Second,
		downloadDelayMax: time.Second * 2,

		searchWorkers:      8,
		searchMaxDistance:  2,
		searchRandomResult: false,
		searchTimeout:      time.Minute * 2,

		providerUserAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0",
		providerTimeout:   time.Second * 30,
	}
}

func GetConfig() *Config {
	defaultCfg := GetDefaultConfig()
	cfg := &Config{}

	cfg.answerCache = stringToBoolFallback(os.Getenv("CACHE_ANSWER"), defaultCfg.answerCache)
	cfg.answerCacheExpiration = stringToDurationFallback(os.Getenv("CACHE_ANSWER_EXPIRATION"), defaultCfg.answerCacheExpiration)
	cfg.answerCacheCleanupInterval = stringToDurationFallback(os.Getenv("CACHE_ANSWER_CLEANUP_INTERVAL"), defaultCfg.answerCacheCleanupInterval)

	cfg.listingCache = stringToBoolFallback(os.Getenv("CACHE_LISTING"), defaultCfg.listingCache)
	cfg.listingCacheExpiration = stringToDurationFallback(os.Getenv("CACHE_LISTING_EXPIRATION"), defaultCfg.listingCacheExpiration)
	cfg.listingCacheCleanupInterval = stringToDurationFallback(os.Getenv("CACHE_LISTING_CLEANUP_INTERVAL"), defaultCfg.listingCacheCleanupInterval)

	cfg.contentCache = stringToBoolFallback(os.Getenv("CACHE_CONTENT"), defaultCfg.contentCache)
	cfg.contentCacheExpiration = stringToDurationFallback(os.Getenv("CACHE_CONTENT_EXPIRATION"), defaultCfg.contentCacheExpiration)
	cfg.contentCacheCleanupInterval = stringToDurationFallback(os.Getenv("CACHE_CONTENT_CLEANUP_INTERVAL"), defaultCfg.contentCacheCleanupInterval)

	cfg.downloadDelayMin = stringToDurationFallback(os.Getenv("DOWNLOAD_DELAY_MIN"), defaultCfg.downloadDelayMin)
	cfg.downloadDelayMax = stringToDurationFallback(os.Getenv("DOWNLOAD_DELAY_MAX"), defaultCfg.downloadDelayMax)

	cfg.searchWorkers = stringToIntFallback(os.Getenv("SEARCH_WORKERS"), defaultCfg.searchWorkers)
	cfg.searchMaxDistance = stringToIntFallback(os.Getenv("SEARCH_MAX_DISTANCE"), defaultCfg.searchMaxDistance)
	cfg.searchRandomResult = stringToBoolFallback(os.Getenv("SEARCH_RANDOM_RESULT"), defaultCfg.searchRandomResult)
	cfg.searchTimeout = stringToDurationFallback(os.Getenv("SEARCH_TIMEOUT"), defaultCfg.searchTimeout)

	cfg.providerUserAgent = stringFallback(os.Getenv("PROVIDER_USER_AGENT"), defaultCfg.providerUserAgent)
	cfg.providerTimeout = stringToDurationFallback(os.Getenv("PROVIDER_TIMEOUT"), defaultCfg.providerTimeout)

	return cfg
}
