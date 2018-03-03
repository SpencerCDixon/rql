package storage

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
)

// func TestNewLogManager(t *testing.T) {
// defer cleanUp("logmanager")
// lm := newLogManager(t, "logmanager")

// testutil.Equals(t, 4, lm.size(10))
// testutil.Equals(t, 9, lm.size("hello"))
// }

func TestAppendingLogManager(t *testing.T) {
	// defer cleanUp("random")
	lm := newLogManager(t, "randomtwo")

	lr1 := []interface{}{"hello", "world"}
	lr2 := []interface{}{1, 2, 3}
	lr3 := []interface{}{42, "meaning", "of", "life"}
	lm.Append(lr1)
	lm.Append(lr2)
	lm.Append(lr3)

	iter := lm.Iterator()

	lm.Flush()
}

func newLogManager(t *testing.T, dbName string) *LogManager {
	t.Helper()
	fm, err := NewFileManager(dbName)
	testutil.Ok(t, err)
	logFile := dbName + ".log"
	lm, err := NewLogManager(logFile, fm)
	testutil.Ok(t, err)
	return lm
}
