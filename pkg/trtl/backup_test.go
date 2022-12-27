package trtl_test

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	honuldb "github.com/rotationalio/honu/engines/leveldb"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/utils"
)

// Test that the backup manager periodically creates backups.
func (s *trtlTestSuite) TestBackupManager() {
	backupDir := "testdata/backup"
	defer os.RemoveAll(backupDir)
	require := s.Require()

	// Restart the trtl service with the backup manager config.
	conf := config.BackupConfig{
		Enabled:  true,
		Interval: time.Millisecond,
		Storage:  backupDir,
		Keep:     1,
	}

	// Create a backup manager that's separate from the trtl service
	backup, err := trtl.NewBackupManager(conf, s.trtl.GetDB())
	require.NoError(err)

	// Start the backup manager
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		backup.Run()
	}()

	// Wait for the backup manager to run through its loop. The shutdown check is at
	// the beginning so there is a timing window here.
	time.Sleep(conf.Interval * 3)

	// Make sure that the backup manager is stopped
	require.NoError(backup.Shutdown())
	wg.Wait()

	// Backup should be created
	require.DirExists(backupDir)
	files, err := os.ReadDir(backupDir)
	require.NoError(err)
	require.Len(files, 1, "wrong number of backups created")
	s.compareBackup(filepath.Join(conf.Storage, files[0].Name()))
}

// Compares the target backup DB to the current DB to verify that they contain the same
// objects.
func (s *trtlTestSuite) compareBackup(path string) {
	require := s.Require()

	// Extract the backup DB
	root, err := utils.ExtractGzip(path, "testdata/lastbackup", false)
	require.NoError(err)
	defer os.RemoveAll(root)

	// Open the backup DB
	backup, err := leveldb.OpenFile(root, nil)
	require.NoError(err)

	// Get the current underlying levelDB object
	engine, ok := s.trtl.GetDB().Engine().(*honuldb.LevelDBEngine)
	require.True(ok)
	current := engine.DB()

	// Make sure everything in the current DB is also in the backup DB
	iter := current.NewIterator(nil, nil)
	for iter.Next() {
		// Make sure the value is the same
		val, err := backup.Get(iter.Key(), nil)
		require.NoError(err)
		require.Equal(val, iter.Value())
	}
	require.NoError(iter.Error())
	iter.Release()

	// Make sure there are no extra keys in the backup DB
	iter = backup.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		_, err := current.Get(key, nil)
		require.NoError(err)
	}
	require.NoError(iter.Error())
	iter.Release()
}
