package config

import (
	"go.rtnl.ai/confire"
	"go.rtnl.ai/x/rlog"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be GDS_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const Prefix = "gds"

// Config contains all of the configuration parameters for the Uptime server which
// are loaded from the environment and should be validated before use.
type Config struct {
	Maintenance    bool              `default:"false" desc:"if true the server will start in maintenance mode"`
	Mode           string            `default:"release" desc:"specify the mode of the gin server (release, debug, testing)"`
	LogLevel       rlog.LevelDecoder `default:"info" split_words:"true" desc:"set the log level for the server"`
	ConsoleLog     bool              `default:"false" split_words:"true" desc:"if true the server will log to the console in text format"`
	BindAddr       string            `default:":8000" split_words:"true" desc:"the ip address and port to bind the server to"`
	AllowedOrigins []string          `split_words:"true" default:"http://localhost:8000" desc:"a list of allowed origins for CORS"`
}

// New creates a new Config instance and loads the configuration from the environment,
// validating the configuration and returning an error if the configuration is invalid
// or could not be parsed from environment variables.
//
// NOTE: New should only be used for testing, for module access to the config use Get().
func New() (conf *Config, err error) {
	// NOTE: confire.Process calls Validate() internally.
	conf = &Config{}
	if err = confire.Process(Prefix, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
