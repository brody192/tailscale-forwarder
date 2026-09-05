package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"main/internal/logger"
	"main/internal/util"

	"github.com/caarlos0/env/v11"
)

var Cfg = config{}

func Load() {
	var errs []error

	if err := env.Parse(&Cfg); err != nil {
		if e, ok := err.(env.AggregateError); ok {
			errs = append(errs, e.Errors...)
		} else {
			errs = append(errs, err)
		}
	}

	connectionMappings, err := parseConnectionMappings("CONNECTION_MAPPING_", os.Environ())
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

	tags := []string{}
	for _, tag := range Cfg.TSAdvertiseTags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if !strings.HasPrefix(tag, "tag:") {
			errs = append(errs, fmt.Errorf("TS_ADVERTISE_TAGS entries must start with %q, got %q", "tag:", tag))
		}
		tags = append(tags, tag)
	}
	Cfg.TSAdvertiseTags = tags

	// tsnet turns an OAuth client secret into an auth key only when tags are
	// advertised; fail early with a readable message instead of tsnet's.
	if strings.HasPrefix(Cfg.TSAuthKey, "tskey-client-") && len(tags) == 0 {
		errs = append(errs, fmt.Errorf("TS_AUTHKEY is an OAuth client secret, which requires TS_ADVERTISE_TAGS (e.g. tag:railway)"))
	}

	if len(errs) > 0 {
		logger.StderrWithSource.Error("configuration error(s) found", logger.ErrorsAttr(errs...))
		os.Exit(1)
	}

	// Tailscale reads TS_DEBUG_MTU directly from the environment when bringing
	// the network up, so write the resolved value back for it to pick up.
	os.Setenv("TS_DEBUG_MTU", strconv.Itoa(Cfg.TSDebugMTU))
}
