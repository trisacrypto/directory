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
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
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
		sentry.Error(nil).Msg("currently configured store cannot be backed up, mark as disabled or use different store")
		return
	}

	// Check backup directory, creating it as necessary
	backupDir, err := s.getBackupStorage()
	if err != nil {
		sentry.Fatal(nil).Err(err).Msg("backup manager cannot access backup directory")
	}

	ticker := time.NewTicker(s.conf.Backup.Interval)
	log.Info().Dur("interval", s.conf.Backup.Interval).Str("store", backupDir).Msg("backup manager started")

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

		// Execute the backup - error messages are logged in the Backup function so they
		// are ignored here and are only returned from Backup for testing purposes.
		s.Backup(backupDir)
	}
}

// Backup performs a backup on behalf of the service. BackupManager calls this function
// at a periodic interval to take a snapshot of the store to disk and to keep only the
// configured number of archives.
func (s *Service) Backup(path string) (err error) {
	// Begin the backup process
	start := time.Now()
	log.Debug().Msg("starting backup")

	// Primarily a testing helper, if the path is empty, attempt to resolve or create
	// the backup directory from the configuration. Passing empty string for the path
	// can also be used for one off backups for utilities that want to call this func.
	if path == "" {
		if path, err = s.getBackupStorage(); err != nil {
			log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not get backup storage")
			return err
		}
	}

	// Conduct the backup, logging errors if needed
	if err = s.db.(store.Backup).Backup(path); err != nil {
		// Do not continue if there was a backup error; all code in the rest of the
		// loop should expect that the backup was successful.
		// NOTE: using WithLevel and Fatal does not Exit the program like log.Fatal()
		// this ensures that we issue a CRITICAL severity without stopping the server.
		log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not backup database")
		return err
	}

	log.Info().Dur("duration", time.Since(start)).Msg("backup complete")

	// Remove any previous backups that may be in the directory
	// NOTE: this requires the backup to write filenames as gdsdb-200601021504.*
	var archives []string
	if archives, err = listArchives(path); err != nil {
		sentry.Error(nil).Err(err).Msg("could not list backup directory")
		return err
	}

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
	return nil
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
