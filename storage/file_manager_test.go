package storage

import (
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spencercdixon/rql/testutil"
)

func TestBlockEquals(t *testing.T) {
	b1 := NewBlock("users.tbl", 0)
	b2 := NewBlock("users.tbl", 1)
	b3 := NewBlock("users.tbl", 0)
	b4 := NewBlock("other.tbl", 0)

	testutil.Assert(t, b1.Equals(b3), "blocks are equal")
	testutil.Assert(t, !b1.Equals(b2), "blocks are not equal")
	testutil.Assert(t, !b1.Equals(b4), "blocks are not equal")
}

func TestBlockString(t *testing.T) {
	b1 := NewBlock("users.tbl", 0)
	testutil.Equals(t, b1.String(), "[file users.tbl, block 0]")
}

func TestFileManager(t *testing.T) {
	defer cleanUp("example")
	fm, err := NewFileManager("example")
	testutil.Ok(t, err)
	p := NewPage(fm)

	// set up blocks in various parts of the file
	blk1 := &Block{FileName: "users.tbl", BlockNum: 0}
	blk2 := &Block{FileName: "users.tbl", BlockNum: 2}

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

func TestPageInt(t *testing.T) {
	defer cleanUp("example")
	fm, err := NewFileManager("example")
	testutil.Ok(t, err)
	p := NewPage(fm)
	p.SetInt(0, 42)
	myInt := p.GetInt(0)
	testutil.Equals(t, 42, myInt)
}

func TestPageString(t *testing.T) {
	defer cleanUp("example")
	fm, err := NewFileManager("example")
	testutil.Ok(t, err)
	p := NewPage(fm)

	p.SetString(0, "hello")
	str := p.GetString(0)
	testutil.Equals(t, "hello", str)
}

func TestPageCombined(t *testing.T) {
	defer cleanUp("example")
	fm, err := NewFileManager("example")
	testutil.Ok(t, err)
	p := NewPage(fm)

	p.SetInt(0, 20)
	p.SetInt(10, 42)
	p.SetString(15, "hello")
	p.SetString(80, "world")

	twenty := p.GetInt(0)
	fourtytwo := p.GetInt(10)
	hello := p.GetString(15)
	world := p.GetString(80)

	testutil.Equals(t, 20, twenty)
	testutil.Equals(t, 42, fourtytwo)
	testutil.Equals(t, "hello", hello)
	testutil.Equals(t, "world", world)
}

// remove the db directories and files that get created while testing
func cleanUp(dbName string) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, "rql", dbName)
	os.RemoveAll(path)
}
