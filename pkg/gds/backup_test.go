package gds_test

import (
	"io/ioutil"
	"os"
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
	go s.svc.BackupManager(nil)

	// Wait for the backup interval to elapse
	time.Sleep(s.svc.GetConf().Backup.Interval * 2)

	// Backup should not be created
	backupDir := s.svc.GetConf().Backup.Storage
	require.NoDirExists(backupDir)
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
	go s.svc.BackupManager(stop)

	// Wait for a the backup interval to elapse
	time.Sleep(s.svc.GetConf().Backup.Interval * 2)

	// Backup should be created
	backupDir := s.svc.GetConf().Backup.Storage
	require.DirExists(backupDir)
	files, err := ioutil.ReadDir(backupDir)
	require.NoError(err)
	numBackups := 0
	for f := range files {
		if files[f].IsDir() {
			numBackups++
		}
	}
	require.Equal(1, numBackups, "wrong number of backups created")
	stop <- true
}
