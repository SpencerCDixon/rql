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
	blk1 := &Block{FileName: "users.tbl", BlockNum: 0}
	p.Read(blk1)
	p.SetString(0, "hello")
	p.Write(blk1)
	str := p.GetString(0)
	testutil.Equals(t, "hello", str)
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
