package config

type connectionMapping struct {
	SourcePort int
	TargetAddr string
	TargetPort int
}

type config struct {
	TSHostname string `env:"TS_HOSTNAME,required"`
	TSAuthKey  string `env:"TS_AUTHKEY,required"`

	ConnectionMappings []connectionMapping
}
