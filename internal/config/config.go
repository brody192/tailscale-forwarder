package config

import (
	"fmt"
	"os"

	"main/internal/logger"
	"main/internal/util"

	"github.com/caarlos0/env/v11"
)

var Cfg = config{}

func init() {
	var errs []error

	if err := env.Parse(&Cfg); err != nil {
		if e, ok := err.(env.AggregateError); ok {
			errs = append(errs, e.Errors...)
		} else {
			errs = append(errs, err)
		}
	}

	connectionMappings, err := parseConnectionMappingsFromEnv("CONNECTION_MAPPING_")
	if err != nil {
		errs = append(errs, err)
	}

	if len(connectionMappings) == 0 && err == nil {
		errs = append(errs, fmt.Errorf("required environment variable \"CONNECTION_MAPPING_[n]\" is not set"))
	}

	Cfg.ConnectionMappings = connectionMappings

	sanitizedHostname := util.SanitizeString(Cfg.TSHostname)

	if sanitizedHostname == "" {
		errs = append(errs, fmt.Errorf("TS_HOSTNAME must be a valid hostname, before sanitization: \"%s\", after sanitization: \"%s\"", Cfg.TSHostname, sanitizedHostname))
	}

	Cfg.TSHostname = sanitizedHostname

	if len(errs) > 0 {
		logger.StderrWithSource.Error("configuration error(s) found", logger.ErrorsAttr(errs...))
		os.Exit(1)
	}
}
