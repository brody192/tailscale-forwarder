package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func parseConnectionMappingsFromEnv(prefix string) ([]connectionMapping, error) {
	connectionMappings := []connectionMapping{}

	for _, envVar := range os.Environ() {
		kv := strings.SplitN(envVar, "=", 2)

		if len(kv) != 2 {
			continue
		}

		if !strings.HasPrefix(kv[0], prefix) {
			continue
		}

		parts := strings.SplitN(kv[1], ":", 4)

		if len(parts) < 3 || len(parts) > 4 {
			return nil, fmt.Errorf("invalid connection mapping: %s (expected <source_port>:<target_host>:<target_port>[:<protocol>])", kv[1])
		}

		sourcePort, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid source port: %s", parts[0])
		}

		targetPort, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid target port: %s", parts[2])
		}

		protocol := ConnectionProtocolHTTP
		if len(parts) == 4 {
			protocol = strings.ToUpper(strings.TrimSpace(parts[3]))
			if protocol != ConnectionProtocolHTTP && protocol != ConnectionProtocolHTTPS {
				return nil, fmt.Errorf("invalid protocol: %s (expected HTTP or HTTPS)", parts[3])
			}
		}

		connectionMappings = append(connectionMappings, connectionMapping{
			SourcePort: sourcePort,
			TargetAddr: parts[1],
			TargetPort: targetPort,
			Protocol:   protocol,
		})
	}

	sourcePorts := []int{}

	for _, connectionMapping := range connectionMappings {
		if slices.Contains(sourcePorts, connectionMapping.SourcePort) {
			return nil, fmt.Errorf("duplicate source port %d found in connection mappings", connectionMapping.SourcePort)
		}

		sourcePorts = append(sourcePorts, connectionMapping.SourcePort)
	}

	return connectionMappings, nil
}
