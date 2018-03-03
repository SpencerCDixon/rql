package storage

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
)

func TestPageInt(t *testing.T) {
	defer cleanUp("example")
	p := newPage(t, "example")

	p.SetInt(0, 42)
	myInt := p.GetInt(0)
	testutil.Equals(t, 42, myInt)
}

func TestPageString(t *testing.T) {
	defer cleanUp("example")
	p := newPage(t, "example")

	p.SetString(0, "hello")
	str := p.GetString(0)
	testutil.Equals(t, "hello", str)
}

func TestPageCombined(t *testing.T) {
	defer cleanUp("example")
	p := newPage(t, "example")

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

func TestSetErrors(t *testing.T) {
	defer cleanUp("example")
	p := newPage(t, "example")

	// not enough room for 4 byte int
	err := p.SetInt(397, 542)
	testutil.Equals(t, err, ErrPageFull)

	// not enough room for string
	err = p.SetString(390, "this is one long string")
	testutil.Equals(t, err, ErrPageFull)
}

func newPage(t *testing.T, dbName string) *Page {
	t.Helper()
	fm, err := NewFileManager(dbName)
	testutil.Ok(t, err)
	return NewPage(fm)
}
