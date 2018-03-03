package storage

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
)

func TestNewLogManager(t *testing.T) {
	defer cleanUp("logmanager")
	lm := newLogManager(t, "logmanager")

	testutil.Equals(t, 4, lm.size(10))
	testutil.Equals(t, 9, lm.size("hello"))
}

func TestAppendingLogManager(t *testing.T) {
	defer cleanUp("logmanager")
	lm := newLogManager(t, "logmanager")

	lr1 := []interface{}{"hello", "world"}
	lr2 := []interface{}{1, 2, 3}
	lr3 := []interface{}{42, "meaning", "of", "life"}
	lm.Append(lr1)
	lm.Append(lr2)
	lm.Append(lr3)

	lm.Flush()
}

func TestLogIterator(t *testing.T) {
	defer cleanUp("iterator")
	lm := newLogManager(t, "iterator")

	lr1 := []interface{}{"hello", "world"}
	lr2 := []interface{}{1, 40}
	lm.Append(lr1)
	lm.Append(lr2)

	// will Flush
	iter := lm.Iterator()

	// Records proper int records
	iter.Next()
	l2 := iter.Value()
	one := l2.NextInt()
	forty := l2.NextInt()

	testutil.Equals(t, 1, one)
	testutil.Equals(t, 40, forty)

	// Records proper string records
	iter.Next()
	l1 := iter.Value()
	hello := l1.NextString()
	world := l1.NextString()
	testutil.Equals(t, hello, "hello")
	testutil.Equals(t, world, "world")
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
