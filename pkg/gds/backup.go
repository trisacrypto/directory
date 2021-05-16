package gds

import "github.com/rs/zerolog/log"

// BackupManager is a go routine that periodically copies the directory storage to a
// compressed backup location, either locally on disk or to a cloud location. The
// manager may also encrypt the storage with provided keys. The manager is started when
// the server is started; but if it is not able to run, it will exit before continuing.
//
// TODO: allow storage to cloud storage rather than to disk
// TODO: encrypt the backup storage file
func (s *Server) BackupManager() {
	if !s.conf.Backup.Enabled {
		log.Warn().Msg("backup manager is not enabled")
		return
	}

	log.Info().Dur("interval", s.conf.Backup.Interval).Str("store", s.conf.Backup.Storage).Msg("backup manager started")
}

// get the configured backup directory storage
