package gds_test

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/fixtures"
)

// Test that the backup manager does not create backups if disabled.
func (s *gdsTestSuite) TestBackupManagerDisabled() {
	conf := gds.MockConfig()
	conf.Backup = config.BackupConfig{
		Enabled:  false,
		Interval: time.Millisecond,
		Storage:  "testdata/backup",
		Keep:     1,
	}
	s.SetConfig(conf)
	defer s.ResetConfig()
	s.LoadEmptyFixtures()
	require := s.Require()

	// Start the backup manager
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.svc.BackupManager(nil)
	}()

	// Backup should not be created
	backupDir := s.svc.GetConf().Backup.Storage
	require.NoDirExists(backupDir)

	// Make sure the backup manager is stopped before we exit
	wg.Wait()
}

// Test that the backup manager periodically creates backups.
func (s *gdsTestSuite) TestBackupManager() {
	if s.fixtures.StoreType() == fixtures.StoreTrtl {
		s.T().Skip("backup manager not supported for trtl store")
	}

	conf := gds.MockConfig()
	conf.Backup = config.BackupConfig{
		Enabled:  true,
		Interval: 100 * time.Millisecond,
		Storage:  "testdata/backups",
		Keep:     1,
	}
	s.SetConfig(conf)
	defer s.ResetConfig()
	s.LoadEmptyFixtures()
	defer os.RemoveAll(s.svc.GetConf().Backup.Storage)
	require := s.Require()

	// Execute a single Backup; testing the looping functionality of the BackupManager
	// will result in a race condition, so we assume that any runtime errors will
	// primarily occur in the Backup function and that the BackupManager routine is ok.
	err := s.svc.Backup("")
	require.NoError(err, "could not execute backup")

	// Backup should be created
	backupDir := s.svc.GetConf().Backup.Storage
	require.DirExists(backupDir)
	files, err := ioutil.ReadDir(backupDir)
	require.NoError(err)
	require.Len(files, 1, "wrong number of backups created")
}
