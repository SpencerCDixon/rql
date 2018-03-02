package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
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
	fm := &FileManager{
		Dir:       dbLoc,
		openFiles: make(map[string]*os.File),
	}

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

func (fm *FileManager) Read(blk *Block, content []byte) error {
	file, err := fm.getFile(blk.FileName)
	if err != nil {
		return err
	}
	offset := blk.BlockNum * BlockSize
	if _, err := file.ReadAt(content, int64(offset)); err != nil {
		return err
	}
	return nil
}

func (fm *FileManager) Write(blk *Block, content []byte) error {
	file, err := fm.getFile(blk.FileName)
	if err != nil {
		return err
	}
	offset := int64(blk.BlockNum * BlockSize)
	if _, err := file.WriteAt(content, offset); err != nil {
		return err
	}
	return nil
}

func (fm *FileManager) Append(filename string, content []byte) (*Block, error) {
	newBlkNum := fm.size(filename)
	blk := &Block{FileName: filename, BlockNum: newBlkNum}

	err := fm.Write(blk, content)
	if err != nil {
		return nil, errors.Wrap(err, "writing content")
	}

	return blk, nil
}

// getFile finds an open descriptor that is being saved or creates a new one
// with the proper settings if not found.
func (fm *FileManager) getFile(filename string) (*os.File, error) {
	if file, ok := fm.openFiles[filename]; ok {
		return file, nil
	}

	path := filepath.Join(fm.Dir, filename)
	// O_SYNC is important because when we write to the file we want to ensure it
	// 100% happens before moving on and there is no delay.
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}

	// Save for later
	fm.openFiles[filename] = file
	return file, nil
}

// size returns the current block number for a given file.
func (fm *FileManager) size(filename string) int {
	file, err := fm.getFile(filename)
	if err != nil {
		return 0
	}

	info, err := file.Stat()
	if err != nil {
		return 0
	}

	bytes := info.Size()
	return int(bytes / BlockSize)
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

// GetInt gets an int at the given offset of this pages contents.
func (p *Page) GetInt(offset int) int {
	end := offset + IntSize
	return int(binary.LittleEndian.Uint32(p.content[offset:end]))
}

// SetInt sets an int at the given offset.
func (p *Page) SetInt(offset int, val int) {
	// TODO: handle int being bigger than room left in page
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, uint32(val))
	addition := buf.Bytes()

	// insert the new bytes
	p.content = append(
		p.content[:offset],
		append(addition, p.content[offset+IntSize:]...)...,
	)
}

// GetString gets a string at the given offset of this pages contents.
func (p *Page) GetString(offset int) string {
	intEnd := offset + IntSize
	numChars := binary.LittleEndian.Uint32(p.content[offset:intEnd])
	strEnd := numChars + uint32(intEnd)
	return string(p.content[intEnd:strEnd])
}

// SetString sets a string at the given offset.
func (p *Page) SetString(offset int, val string) {
	// TODO: handle string being bigger than room left in page
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

func (p *Page) Read(blk *Block) {
	p.fm.Read(blk, p.content)
}
func (p *Page) Write(blk *Block) {
	p.fm.Write(blk, p.content)
}
func (p *Page) Append(filename string) (*Block, error) {
	return p.fm.Append(filename, p.content)
}

// reset wipes the contents clean but preserves the underlying storage for use
// by future writes.
func (p *Page) reset() {
	for i := range p.content {
		p.content[i] = 0
	}
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
