package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/guregu/null.v4"
)

type Config struct {
	Fields []string `json:"fileds" envconfig:"K6_OUTPUTERROR_FIELDS"`
	FTime  null.Int `json:"ftime"  envconfig:"K6_OUTPUTERROR_FTIME"`
}

// Create default config
func newConfig() Config {
	return Config{
		FTime: null.NewInt(1000, false),
	}
}

func GetConsolidatedConfig(jsonRawCfg json.RawMessage, confArg string) (Config, error) {
	consCfg := newConfig()

	if jsonRawCfg != nil {
		jsonCfg, err := parseJSON(jsonRawCfg)
		if err != nil {
			return consCfg, err
		}
		consCfg = consCfg.apply(jsonCfg)
	}

	envCfg, err := parseEnv("K6_OUTPUTERROR_FIELDS")
	if err != nil {
		return consCfg, err
	}
	consCfg = consCfg.apply(envCfg)

	if confArg != "" {
		argCfg, err := parseFieldsArg(confArg)
		if err != nil {
			return consCfg, err
		}
		consCfg = consCfg.apply(argCfg)
	}

	for k, v := range consCfg.Fields {
		consCfg.Fields[k] = strings.TrimSpace(v)
	}

	return consCfg, nil
}

func (c Config) apply(cfg Config) Config {
	if len(cfg.Fields) > 0 {
		c.Fields = cfg.Fields
	}
	if cfg.FTime.Valid && cfg.FTime.Int64 > 0 {
		c.FTime = cfg.FTime
	}
	return c
}

func parseEnv(prefix string) (Config, error) {
	cfg := Config{}
	if fields, exist := os.LookupEnv(prefix); exist {
		cfg.Fields = strings.Split(fields, ",")
		if len(cfg.Fields) < 1 {
			return cfg, fmt.Errorf("could not parse env fields")
		}
	}
	return cfg, nil
}

func parseJSON(data json.RawMessage) (Config, error) {
	cfg := Config{}
	err := json.Unmarshal(data, &cfg)
	return cfg, err
}

func parseFieldsArg(arg string) (Config, error) {
	cfg := Config{}
	fields := strings.Split(strings.TrimPrefix(arg, "fields="), ",")
	if len(fields) < 1 {
		return cfg, fmt.Errorf("could not parse arguments fields")
	}
	cfg.Fields = append(cfg.Fields, fields...)

	return cfg, nil
}
