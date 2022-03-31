package gds

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/store"
)

// BackupManager is a go routine that periodically copies the directory storage to a
// compressed backup location, either locally on disk or to a cloud location. The
// manager may also encrypt the storage with provided keys. The manager is started when
// the server is started; but if it is not able to run, it will exit before continuing.
func (s *Service) BackupManager(stop <-chan bool) {
	if !s.conf.Backup.Enabled {
		log.Warn().Msg("backup manager is not enabled")
		return
	}

	// Test that the store is backupable
	if _, ok := s.db.(store.Backup); !ok {
		log.Error().Msg("currently configured store cannot be backed up, mark as disabled or use different store")
		return
	}

	// Check backup directory
	backupDir, err := s.getBackupStorage()
	if err != nil {
		log.Fatal().Err(err).Msg("backup manager cannot access backup directory")
	}

	ticker := time.NewTicker(s.conf.Backup.Interval)
	log.Info().Dur("interval", s.conf.Backup.Interval).Str("store", backupDir).Msg("backup manager started")

backups:
	for {
		// Wait for next tick or a stop message
		select {
		case done := <-stop:
			// The value of the signal doesn't matter, but we check it here for completeness
			if done {
				log.Warn().Msg("backup manager received stop signal")
				return
			}
		case <-ticker.C:
		}

		// Begin the backup process
		start := time.Now()
		log.Debug().Msg("starting backup")

		// Conduct the backup, logging errors if needed
		if err := s.db.(store.Backup).Backup(backupDir); err != nil {
			// Do not continue if there was a backup error; all code in the rest of the
			// loop should expect that the backup was successful.
			// NOTE: using WithLevel and Fatal does not Exit the program like log.Fatal()
			// this ensures that we issue a CRITICAL severity without stopping the server.
			log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not backup database")
			continue backups
		} else {
			log.Info().Dur("duration", time.Since(start)).Msg("backup complete")
		}

		// Remove any previous backups that may be in the directory
		// NOTE: this requires the backup to write filenames as gdsdb-200601021504.*
		var archives []string
		if archives, err = listArchives(backupDir); err != nil {
			log.Error().Err(err).Msg("could not list backup directory")
		} else {
			if len(archives) > s.conf.Backup.Keep {
				var removed int
				for _, archive := range archives[:len(archives)-s.conf.Backup.Keep] {
					log.Debug().Str("archive", archive).Msg("deleting archive")
					if err = os.Remove(archive); err == nil {
						removed++
					}
				}
				log.Debug().Int("kept", s.conf.Backup.Keep).Int("removed", removed).Msg("backup directory cleaned up")
			}
		}
	}
}

// get the configured backup directory storage or return an error
func (s *Service) getBackupStorage() (path string, err error) {
	if s.conf.Backup.Storage == "" {
		return "", errors.New("incorrectly configured: backups enabled but no backup storage")
	}

	var stat os.FileInfo
	path = s.conf.Backup.Storage
	if stat, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// Create the directory if it does not exist
			if err = os.MkdirAll(path, 0755); err != nil {
				return "", fmt.Errorf("could not create backup storage directory: %s", err)
			}
			return path, nil
		}
		return "", err
	}

	if !stat.IsDir() {
		return "", errors.New("incorrectly configured: backup storage is not a directory")
	}
	return path, nil
}

// list all backup archives ordered by date ascending using string sorting that depends
// on the backup archive format gdsdb-YYYYmmddHHMM.
func listArchives(path string) (paths []string, err error) {
	if paths, err = filepath.Glob(filepath.Join(path, "gdsdb-[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9].*")); err != nil {
		return nil, err
	}

	// Sort the paths by timestamp ascending
	sort.Strings(paths)
	return paths, nil
}
