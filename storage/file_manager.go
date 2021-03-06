package storage

import (
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

// NewFileManager returns a new file manager.  It checks for the existence of
// the database 'db' and if none exists it will create the necessary
// directories/files. TODO: tmp file removal.
func NewFileManager(db string) (*FileManager, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	// Store all databases in ~/rql/dbname
	dbLoc := filepath.Join(home, "rql", db)
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

// Write writes the content bytes to the given blk.  If the file does not exist
// for the block a new one will be created and saved in the managers pool.
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

// Append appends content to the filename returning the Block that the bytes
// were written to.
func (fm *FileManager) Append(filename string, content []byte) (*Block, error) {
	newBlkNum, err := fm.Size(filename)
	if err != nil {
		return nil, err
	}

	blk := NewBlock(filename, newBlkNum)

	err = fm.Write(blk, content)
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

// Size returns the current block number for a given file.
func (fm *FileManager) Size(filename string) (int, error) {
	file, err := fm.getFile(filename)
	if err != nil {
		return 0, err
	}

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	bytes := info.Size()
	return int(bytes / BlockSize), nil
}
