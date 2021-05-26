package shared

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB is an instance of the LevelDB client for managing persistent cache.
var LevelDB *leveldb.DB

func initLevelDB() {
	var (
		path = viper.GetString("persistence_path")
		err error
	)

	if len(path) == 0 {
		Logger.Fatal("failed to initialise LevelDB: local path not provided")
	}

	if LevelDB, err = leveldb.OpenFile(path, nil); err != nil {
		Logger.Fatal(errors.Wrap(err, "failed to initialise LevelDB"))
	}
}

func closeLevelDB() {
	if err := LevelDB.Close(); err != nil {
		Logger.Error(errors.Wrap(err, "failed to close connection to LevelDB"))
	}
}
