package mock

import (
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

const (
	metaRegion = "tauceti"
	metaOwner  = "taurian"
	metaPID    = 8
)

// Config returns a configuration that ensures the server will operate in a fully
// mocked way with all testing parameters set correctly. The config is returned directly
// for required modifications, such as pointing the database path to a fixtures path.
func Config() config.Config {
	return config.Config{
		Maintenance: false,
		BindAddr:    ":4436",
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  false,
		Database: config.DatabaseConfig{
			URL:           "leveldb:///testdata/db",
			ReindexOnBoot: false,
		},
		Replica: config.ReplicaConfig{
			Enabled: false,
			PID:     metaPID,
			Region:  metaRegion,
			Name:    metaOwner,
		},
		MTLS: config.MTLSConfig{
			Insecure: true,
		},
		Backup: config.BackupConfig{
			Enabled: false,
		},
	}
}
