package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// FileManager is the component responsible for interacting with the operating
// system.
type FileManager struct {
	// Dir is the location of the directory for our database
	Dir string
	// IsNew is a flag to represent whether or not the file manager created the
	// new directory for this database to live in.
	IsNew bool
	// openFiles are all the files that have been opened and are currently in use
	openFiles map[string]*os.File
}

// Block is unit of measurement for how pages manage disk access.
type Block struct {
	FileName string
	BlockNum int
}

// Page
type Page struct {
	content []byte
	fm      *FileManager
}

//-------------
// File Manager
//-------------
func NewFileManager(dir string) (*FileManager, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	// Store all databases in ~/rql/dbname
	dbLoc := filepath.Join(home, "rql", dir)

	fm := &FileManager{Dir: dbLoc}
	if ok := exists(dbLoc); !ok {
		fm.IsNew = true
		if err := os.MkdirAll(dbLoc, 0777); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	// TODO: remove any 'temp' files in the database from previous boots

	return fm, nil
}

//-------------
// Block
//-------------
func NewBlock(filename string, blockNum int) *Block {
	return &Block{FileName: filename, BlockNum: blockNum}
}

// Equals let's us determine if two blocks are the same or not.  If they are
// from the same file and located at the same offset then they are the same
// block.
func (b *Block) Equals(other *Block) bool {
	return other.FileName == b.FileName && other.BlockNum == b.BlockNum
}

// String lets us pretty print blocks for debugging purposes: [file users.tbl, block 3]
func (b *Block) String() string {
	return fmt.Sprintf("[file %s, block %d]", b.FileName, b.BlockNum)
}

//-------------
// Page
//-------------
const (
	BlockSize = 400
	// IntSize represents how many bytes we will let the INT type be in our
	// database. We're using uint32 so this will be 4 bytes.
	IntSize = 4
)

func NewPage(fm *FileManager) *Page {
	content := make([]byte, BlockSize, BlockSize)

	return &Page{
		content: content,
		fm:      fm,
	}
}

func (p *Page) GetInt(offset int) int {
	end := offset + IntSize
	return int(binary.LittleEndian.Uint32(p.content[offset:end]))
}

func (p *Page) SetInt(offset int, val int) {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, uint32(val))
	addition := buf.Bytes()

	// insert the new bytes
	p.content = append(
		p.content[:offset],
		append(addition, p.content[offset+IntSize:]...)...,
	)
}

func (p *Page) GetString(offset int) string {
	intEnd := offset + IntSize
	numChars := binary.LittleEndian.Uint32(p.content[offset:intEnd])
	strEnd := numChars + uint32(intEnd)
	return string(p.content[intEnd:strEnd])
}

func (p *Page) SetString(offset int, val string) {
	buf := bytes.NewBuffer(nil)
	strSizeInt := stringSize(val)
	binary.Write(buf, binary.LittleEndian, uint32(strSizeInt))
	binary.Write(buf, binary.LittleEndian, []byte(val))

	addition := buf.Bytes()
	// insert the new bytes
	p.content = append(
		p.content[:offset],
		append(addition, p.content[offset+IntSize+strSizeInt:]...)...,
	)
}

//-------------
// Utility
//-------------

// Exists returns a bool of wether or not a path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// stringSize returns the number of bytes to allocate for this string
func stringSize(s string) int {
	return len(s)
}
