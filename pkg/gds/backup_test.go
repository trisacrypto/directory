package gds_test

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
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
	conf := gds.MockConfig()
	conf.Backup = config.BackupConfig{
		Enabled:  true,
		Interval: time.Millisecond,
		Storage:  "testdata/backups",
		Keep:     1,
	}
	s.SetConfig(conf)
	defer s.ResetConfig()
	s.LoadEmptyFixtures()
	defer os.RemoveAll(s.svc.GetConf().Backup.Storage)
	require := s.Require()

	// Start the backup manager
	stop := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.svc.BackupManager(stop)
	}()

	// Wait for the backup manager to run through its loop. The shutdown check is at
	// the beginning so there is a timing window here.
	time.Sleep(s.svc.GetConf().Backup.Interval * 2)

	// Make sure that the backup manager is stopped before we proceed
	stop <- true
	wg.Wait()

	// Backup should be created
	backupDir := s.svc.GetConf().Backup.Storage
	require.DirExists(backupDir)
	files, err := ioutil.ReadDir(backupDir)
	require.NoError(err)
	require.Len(files, 1, "wrong number of backups created")
}
