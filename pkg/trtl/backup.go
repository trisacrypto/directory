package trtl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/rotationalio/honu"
	honuldb "github.com/rotationalio/honu/engines/leveldb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	ldbstore "github.com/trisacrypto/directory/pkg/gds/store/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/utils"
)

// BackupManager runs as an independent service which periodically copies the trtl
// storage to a compressed backup location on disk.
type BackupManager struct {
	conf config.BackupConfig
	db   *honu.DB
	stop chan struct{}
}

func NewBackupManager(s *Server) (*BackupManager, error) {
	return &BackupManager{
		conf: s.conf.Backup,
		db:   s.db,
	}, nil
}

// Runs the main BackupManager routine which periodically wakes up and creates a backup
// of the trtl database.
func (m *BackupManager) Run() {
	if !m.conf.Enabled {
		log.Warn().Msg("trtl backups disabled")
		return
	}

	// Create the shutdown channel in Run so that calls to Shutdown() then Run() work
	// NOTE: no buffer in the channel to ensure that the Shutdown caller blocks until
	// the backup has been completed.
	m.stop = make(chan struct{})

	backupDir, err := m.getBackupStorage()
	if err != nil {
		// If we're here, we're just starting the BackupManager, do not continue if we
		// know we won't be able to access the backup directory. log.Fatal() will kill
		// the program with an exit code of 1.
		log.Fatal().Err(err).Msg("trtl backup manager cannot access backup directory")
	}

	ticker := time.NewTicker(m.conf.Interval)
	log.Info().Dur("interval", m.conf.Interval).Str("store", backupDir).Msg("trtl backup manager started")

backups:
	for {
		// Wait for next tick or a stop message
		select {
		case <-m.stop:
			log.Warn().Msg("trtl backup manager received stop signal")
			return
		case <-ticker.C:
		}

		// Begin the backup process
		start := time.Now()
		log.Debug().Msg("starting backup of trtl database")

		// Perform the backup
		if err = m.backup(backupDir); err != nil {
			// Do not continue if there was a backup error; all code in the rest of the
			// loop should expect that the backup was successful.
			// NOTE: using WithLevel and Fatal does not Exit the program like log.Fatal()
			// this ensures that we issue a CRITICAL severity without stopping the server.
			log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not backup database")
			continue backups
		} else {
			log.Info().Dur("duration", time.Since(start)).Msg("trtl backup complete")
		}

		// Remove any previous backups that may be in the directory
		// NOTE: this requires the backup to write filenames as trtldb-200601021504.*
		var archives []string
		if archives, err = listArchives(backupDir); err != nil {
			log.Error().Err(err).Msg("could not list backup directory")
		} else {
			if len(archives) > m.conf.Keep {
				var removed int
				for _, archive := range archives[:len(archives)-m.conf.Keep] {
					log.Debug().Str("archive", archive).Msg("deleting archive")
					if err = os.Remove(archive); err == nil {
						removed++
					}
				}
				log.Debug().Int("kept", m.conf.Keep).Int("removed", removed).Msg("backup directory cleaned up")
			}
		}
	}
}

func (m *BackupManager) Shutdown() error {
	if m.stop != nil {
		// Should block until the current backup completes
		m.stop <- struct{}{}

		// Close the channel and set it to nil so that multiple shutdown calls don't block.
		close(m.stop)
		m.stop = nil
	}
	return nil
}

func (m *BackupManager) backup(path string) (err error) {
	// Create the directory for the copied honu database
	archive := filepath.Join(path, time.Now().UTC().Format("trtldb-200601021504"))
	if err = os.Mkdir(archive, 0755); err != nil {
		return fmt.Errorf("could not create archive directory: %s", err)
	}

	// Ensure the archive directory is cleaned up when the backup is complete
	defer func() {
		os.RemoveAll(archive)
	}()

	// Open a second leveldb database at the backup location
	arcdb, err := leveldb.OpenFile(archive, nil)
	if err != nil {
		return fmt.Errorf("could not open archive database: %s", err)
	}

	// Get the underlying levelDB object
	var engine *honuldb.LevelDBEngine
	var ok bool
	if engine, ok = m.db.Engine().(*honuldb.LevelDBEngine); !ok {
		return fmt.Errorf("unexpected database engine type: %T, expected %T", engine, &honuldb.LevelDBEngine{})
	}
	ldb := engine.DB()

	// Copy all records to the archive database
	var narchived uint64
	if narchived, err = ldbstore.CopyDB(ldb, arcdb); err != nil {
		return fmt.Errorf("could not write all records to archive database, wrote %d records: %s", narchived, err)
	}
	log.Info().Uint64("records", narchived).Msg("trtl archive completed")

	// Close the archive database
	if err = arcdb.Close(); err != nil {
		return fmt.Errorf("could not close archive database: %s", err)
	}

	// Create the compressed tar archive
	dest := filepath.Join(filepath.Dir(archive), filepath.Base(archive)+".tgz")
	if err = utils.WriteGzip(archive, dest); err != nil {
		return fmt.Errorf("could not create compressed tar archive: %s", err)
	}

	// Remove the archive database
	if err = os.RemoveAll(archive); err != nil {
		return fmt.Errorf("could not remove archive database: %s", err)
	}

	return nil
}

// get the configured backup directory storage or return an error
func (m *BackupManager) getBackupStorage() (path string, err error) {
	if m.conf.Storage == "" {
		return "", errors.New("incorrectly configured: backups enabled but no backup storage")
	}

	var stat os.FileInfo
	path = m.conf.Storage
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
// on the backup archive format trtldb-YYYYmmddHHMM.
func listArchives(path string) (paths []string, err error) {
	if paths, err = filepath.Glob(filepath.Join(path, "trtldb-[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9].*")); err != nil {
		return nil, err
	}

	// Sort the paths by timestamp ascending
	sort.Strings(paths)
	return paths, nil
}
