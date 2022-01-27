package trtl_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/options"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/utils"
	"google.golang.org/protobuf/proto"
)

// Test that the backup manager does not create backups if disabled.
func (s *trtlTestSuite) TestBackupManagerDisabled() {
	defer s.reset()
	require := s.Require()

	s.trtl.Shutdown()

	s.conf.Backup = config.BackupConfig{
		Enabled:  false,
		Interval: time.Millisecond,
		Storage:  "testdata/backup",
		Keep:     1,
	}
	var err error
	s.trtl, err = trtl.New(*s.conf)
	require.NoError(err)

	backup, err := trtl.NewBackupManager(s.trtl)
	require.NoError(err)

	// Start the backup manager
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		go backup.Run(nil)
	}()

	// Backup should not be created
	backupDir := s.conf.Backup.Storage
	require.NoDirExists(backupDir)

	// Make sure the backup manager is stopped before we exit
	wg.Wait()
}

// Test that the backup manager periodically creates backups.
func (s *trtlTestSuite) TestBackupManager() {
	defer s.reset()
	defer os.RemoveAll(s.conf.Backup.Storage)
	require := s.Require()

	s.trtl.Shutdown()

	s.conf.Backup = config.BackupConfig{
		Enabled:  true,
		Interval: time.Millisecond,
		Storage:  "testdata/backups",
		Keep:     1,
	}
	var err error
	s.trtl, err = trtl.New(*s.conf)
	require.NoError(err)

	backup, err := trtl.NewBackupManager(s.trtl)
	require.NoError(err)

	// Start the backup manager
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		backup.Run(stop)
	}()

	// Wait for at least one backup interval to elapse
	time.Sleep(s.conf.Backup.Interval)

	// Make sure that the backup manager is stopped before we proceed
	backup.Shutdown()
	wg.Wait()

	// Backup should be created
	backupDir := s.conf.Backup.Storage
	require.DirExists(backupDir)
	files, err := ioutil.ReadDir(backupDir)
	require.NoError(err)
	require.Len(files, 1, "wrong number of backups created")
	s.compareBackup(files[0].Name())
}

// Compares the target backup DB to the current DB to verify that they contain the same
// objects.
func (s *trtlTestSuite) compareBackup(name string) {
	require := s.Require()

	// Extract the backup DB
	root, err := utils.ExtractGzip(filepath.Join(s.conf.Backup.Storage, name), "testdata/lastbackup", false)
	require.NoError(err)
	defer os.RemoveAll(root)

	// Open the backup DB
	backup, err := honu.Open("leveldb:///"+root, s.conf.GetHonuConfig())
	require.NoError(err)

	// Make sure both databases have the same objects
	iter, err := s.trtl.GetDB().Iter(nil)
	require.NoError(err)
	for iter.Next() {
		dbObject, err := iter.Object()
		require.NoError(err)

		backupObject, err := backup.Object(dbObject.Key, options.WithNamespace(dbObject.Namespace))
		require.NoError(err)

		require.True(proto.Equal(dbObject, backupObject), "objects do not match")
	}
	iter.Release()
	require.NoError(iter.Error())
}
