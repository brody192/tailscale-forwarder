package config

type connectionMapping struct {
	SourcePort int
	TargetAddr string
	TargetPort int
	HTTPS      bool
}

type config struct {
	TSHostname   string `env:"TS_HOSTNAME,required"`
	TSAuthKey    string `env:"TS_AUTHKEY,required"`
	TSControlURL string `env:"TS_CONTROL_URL"`
	TSStateDir   string `env:"TS_STATE_DIR"`
	TSEphemeral  bool   `env:"TS_EPHEMERAL" envDefault:"true"`

	// Tags advertised by the node, e.g. "tag:railway,tag:proxy". Required when
	// TS_AUTHKEY is an OAuth client secret (tskey-client-…): tsnet exchanges the
	// secret for an auth key carrying exactly these tags. Optional otherwise.
	TSAdvertiseTags []string `env:"TS_ADVERTISE_TAGS" envSeparator:","`

	// Tailscale's default tun MTU (1280) plus WireGuard overhead exceeds some
	// network paths' MTU, silently stalling large transfers via PMTU
	// blackholing while small ones succeed. Default to a conservative value
	// that fits within common constrained paths.
	TSDebugMTU int `env:"TS_DEBUG_MTU" envDefault:"1230"`

	ConnectionMappings []connectionMapping
}
