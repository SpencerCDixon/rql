package storage

import (
	"bytes"
	"encoding/binary"
)

// Pages are used by the file manager to read and write blocks of bytes.
type Page struct {
	content []byte
	fm      *FileManager
}

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
func (p *Page) SetInt(offset int, val int) error {
	if BlockSize-offset < IntSize {
		return ErrPageFull
	}

	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, uint32(val))
	newBytes := buf.Bytes()

	// insert the new bytes
	p.content = append(
		p.content[:offset],
		append(newBytes, p.content[offset+IntSize:]...)...,
	)

	return nil
}

// GetString gets a string at the given offset of this pages contents.
func (p *Page) GetString(offset int) string {
	intEnd := offset + IntSize
	numChars := binary.LittleEndian.Uint32(p.content[offset:intEnd])
	strEnd := numChars + uint32(intEnd)
	return string(p.content[intEnd:strEnd])
}

// SetString sets a string at the given offset.
func (p *Page) SetString(offset int, val string) error {
	strSizeInt := stringSize(val)

	if BlockSize-offset < IntSize+strSizeInt {
		return ErrPageFull
	}

	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, uint32(strSizeInt))
	binary.Write(buf, binary.LittleEndian, []byte(val))

	newBytes := buf.Bytes()
	// insert the new bytes
	p.content = append(
		p.content[:offset],
		append(newBytes, p.content[offset+IntSize+strSizeInt:]...)...,
	)

	return nil
}

// Read resets our byte slice and then reads the correct block offset of a file
// into the byte slice to be used for setting/getting and writing.
func (p *Page) Read(blk *Block) {
	p.reset()
	p.fm.Read(blk, p.content)
}

// Write persists the pages contents to disk in a synchronous manner.
func (p *Page) Write(blk *Block) {
	p.fm.Write(blk, p.content)
}

// Append increments to the next available block and appends the bytes in this
// pages contents.
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
