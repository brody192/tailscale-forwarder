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

		parts := strings.SplitN(kv[1], ":", 3)

		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid connection mapping: %s", kv[1])
		}

		sourcePort, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid source port: %s", parts[0])
		}

		targetPort, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid target port: %s", parts[2])
		}

		connectionMappings = append(connectionMappings, connectionMapping{
			SourcePort: sourcePort,
			TargetAddr: parts[1],
			TargetPort: targetPort,
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
