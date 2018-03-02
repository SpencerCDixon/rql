package storage

import (
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spencercdixon/rql/testutil"
)

func TestFileManager(t *testing.T) {
	defer cleanUp("example")
	fm, err := NewFileManager("example")
	testutil.Ok(t, err)
	p := NewPage(fm)

	// set up blocks in various parts of the file
	blk1 := NewBlock("users.tbl", 0)
	blk2 := NewBlock("users.tbl", 2)

	// read/write string and int
	p.Read(blk1)
	p.SetString(0, "hello")
	p.SetInt(250, 42)
	p.Write(blk1)
	hello := p.GetString(0)
	life := p.GetInt(250)
	testutil.Equals(t, "hello", hello)
	testutil.Equals(t, 42, life)

	// read/write second block
	p.Read(blk2)
	p.SetString(0, "hello")
	p.SetInt(100, 42)
	p.Write(blk2)
	hello = p.GetString(0)
	life = p.GetInt(100)
	testutil.Equals(t, "hello", hello)
	testutil.Equals(t, 42, life)

	// confirm our first block still persisted
	p.Read(blk1)
	hello = p.GetString(0)
	testutil.Equals(t, "hello", hello)

	// confirm our second block persisted
	p.Read(blk2)
	life = p.GetInt(100)
	testutil.Equals(t, 42, life)
}

// remove the db directories and files that get created while testing
func cleanUp(dbName string) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, "rql", dbName)
	os.RemoveAll(path)
}
