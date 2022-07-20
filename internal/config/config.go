package config

import (
	"fmt"
	"github.com/navcoin/navexplorer-indexer-go/v2/internal/log"
	"github.com/getsentry/raven-go"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Logging            bool
	LogPath            string
	Network            string
	Index              string
	Debug              bool
	Reindex            bool
	ReindexSize        uint64
	RewindToHeight     uint64
	BulkIndex          bool
	BulkTargetHeight   uint64
	BulkIndexSize      uint64
	Subscribe          bool
	SoftForkBlockCycle int
	SoftForkQuorum     int
	Navcoind           NavcoindConfig
	ElasticSearch      ElasticSearchConfig
	ZeroMq             ZeroMqConfig
	Sentry             SentryConfig
}

type NavcoindConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Ssl      bool
	Timeout  int
}

type ElasticSearchConfig struct {
	Hosts       []string
	Sniff       bool
	HealthCheck bool
	Debug       bool
	Username    string
	Password    string
	MappingDir  string
}

type ZeroMqConfig struct {
	Address string
}

type SentryConfig struct {
	Active bool
	DSN    string
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("Unable to init config")
	}

	if Get().Sentry.Active {
		_ = raven.SetDSN(Get().Sentry.DSN)
	}

	initLogger()
}

func initLogger() {
	log.NewLogger(fmt.Sprintf("%s/indexer.log", Get().LogPath), Get().Debug)
}

func Get() *Config {
	return &Config{
		Logging:            getBool("LOGGING", false),
		LogPath:            getString("LOG_PATH", "/app/logs"),
		Network:            getString("NAVCOIND_NETWORK", "mainnet"),
		Index:              getString("INDEX_NAME", "xxx"),
		SoftForkBlockCycle: getInt("SOFTFORK_BLOCKCYCLE", 20160),
		SoftForkQuorum:     getInt("SOFTFORK_QUORUM", 75),
		Debug:              getBool("DEBUG", false),
		Reindex:            getBool("REINDEX", false),
		ReindexSize:        getUint64("REINDEX_SIZE", 200),
		RewindToHeight:     getUint64("REWIND_TO_HEIGHT", 0),
		BulkIndex:          getBool("BULK_INDEX", false),
		BulkTargetHeight:   getUint64("BULK_TARGET_HEIGHT", 0),
		BulkIndexSize:      getUint64("BULK_INDEX_SIZE", 100),
		Subscribe:          getBool("SUBSCRIBE", true),
		Navcoind: NavcoindConfig{
			Host:     getString("NAVCOIND_HOST", ""),
			Port:     getInt("NAVCOIND_PORT", 8332),
			User:     getString("NAVCOIND_USER", "user"),
			Password: getString("NAVCOIND_PASSWORD", "password"),
			Ssl:      getBool("NAVCOIND_SSL", false),
			Timeout:  getInt("NAVCOIND_TIMEOUT", 30),
		},
		ElasticSearch: ElasticSearchConfig{
			Hosts:       getSlice("ELASTIC_SEARCH_HOSTS", make([]string, 0), ","),
			Sniff:       getBool("ELASTIC_SEARCH_SNIFF", true),
			HealthCheck: getBool("ELASTIC_SEARCH_HEALTH_CHECK", true),
			Debug:       getBool("ELASTIC_SEARCH_DEBUG", false),
			Username:    getString("ELASTIC_SEARCH_USERNAME", ""),
			Password:    getString("ELASTIC_SEARCH_PASSWORD", ""),
			MappingDir:  getString("ELASTIC_SEARCH_MAPPING_DIR", "/data/mappings"),
		},
		ZeroMq: ZeroMqConfig{
			Address: getString("ZEROMQ_ADDRESS", "tcp://navcoind:28332"),
		},
		Sentry: SentryConfig{
			Active: getBool("SENTRY_ACTIVE", false),
			DSN:    getString("SENTRY_DSN", ""),
		},
	}
}

func getString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getInt(key string, defaultValue int) int {
	valStr := getString(key, "")
	val, _, err := big.ParseFloat(valStr, 10, 0, big.ToNearestEven)
	if err != nil {
		return defaultValue
	}

	intVal, _ := val.Int64()
	return int(intVal)
}

func getUint(key string, defaultValue uint) uint {
	return uint(getInt(key, int(defaultValue)))
}

func getUint64(key string, defaultValue uint) uint64 {
	return uint64(getInt(key, int(defaultValue)))
}

func getBool(key string, defaultValue bool) bool {
	valStr := getString(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultValue
}

func getSlice(key string, defaultVal []string, sep string) []string {
	valStr := getString(key, "")
	if valStr == "" {
		return defaultVal
	}

	return strings.Split(valStr, sep)
}
